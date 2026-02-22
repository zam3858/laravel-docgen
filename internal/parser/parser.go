package parser

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/zam3858/laravel-docgen/internal/cache"
	"github.com/zam3858/laravel-docgen/internal/model"
)

type Options struct {
	Incremental bool
	Cache       *cache.Cache
}

type Parser struct {
	opts Options
}

func NewParser(opts Options) *Parser {
	return &Parser{opts: opts}
}

func (p *Parser) Parse(root string) (model.Project, []string, error) {
	var project model.Project
	var warnings []string
	var warningsMu sync.Mutex

	addWarnings := func(in []string) {
		if len(in) == 0 {
			return
		}
		warningsMu.Lock()
		defer warningsMu.Unlock()
		warnings = append(warnings, in...)
	}

	filesByKind := map[string][]string{
		"routes":      {},
		"controllers": {},
		"models":      {},
		"services":    {},
		"middleware":  {},
		"jobs":        {},
		"events":      {},
		"listeners":   {},
		"policies":    {},
	}

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".php") {
			return nil
		}

		kind := classify(path, root)
		if kind == "" {
			return nil
		}

		if p.opts.Incremental && p.opts.Cache != nil {
			st, statErr := d.Info()
			if statErr != nil {
				addWarnings([]string{fmt.Sprintf("stat %s: %v", path, statErr)})
				return nil
			}
			if p.opts.Cache.Unchanged(path, st.ModTime()) {
				return nil
			}
			p.opts.Cache.Touch(path, st.ModTime())
		}

		filesByKind[kind] = append(filesByKind[kind], path)
		return nil
	})
	if err != nil {
		return project, warnings, fmt.Errorf("walking project files: %w", err)
	}

	for k := range filesByKind {
		sort.Strings(filesByKind[k])
	}

	var wg sync.WaitGroup
	errCh := make(chan error, 9)

	run := func(fn func() error) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := fn(); err != nil {
				errCh <- err
			}
		}()
	}

	run(func() error {
		result, warn, err := parseRoutes(filesByKind["routes"])
		if err != nil {
			return err
		}
		project.Routes = result
		addWarnings(warn)
		return nil
	})
	run(func() error {
		result, warn, err := parseControllers(filesByKind["controllers"])
		if err != nil {
			return err
		}
		project.Controllers = result
		addWarnings(warn)
		return nil
	})
	run(func() error {
		result, warn, err := parseModels(filesByKind["models"])
		if err != nil {
			return err
		}
		project.Models = result
		addWarnings(warn)
		return nil
	})
	run(func() error {
		result, warn, err := parseServices(filesByKind["services"])
		if err != nil {
			return err
		}
		project.Services = result
		addWarnings(warn)
		return nil
	})
	run(func() error {
		result, warn, err := parseMiddleware(filesByKind["middleware"])
		if err != nil {
			return err
		}
		project.Middleware = result
		addWarnings(warn)
		return nil
	})
	run(func() error {
		result, warn, err := parseJobs(filesByKind["jobs"])
		if err != nil {
			return err
		}
		project.Jobs = result
		addWarnings(warn)
		return nil
	})
	run(func() error {
		result, warn, err := parseEvents(filesByKind["events"])
		if err != nil {
			return err
		}
		project.Events = result
		addWarnings(warn)
		return nil
	})
	run(func() error {
		result, warn, err := parseListeners(filesByKind["listeners"])
		if err != nil {
			return err
		}
		project.Listeners = result
		addWarnings(warn)
		return nil
	})
	run(func() error {
		result, warn, err := parsePolicies(filesByKind["policies"])
		if err != nil {
			return err
		}
		project.Policies = result
		addWarnings(warn)
		return nil
	})

	wg.Wait()
	close(errCh)
	for err := range errCh {
		if err != nil {
			return project, warnings, fmt.Errorf("parsing laravel project: %w", err)
		}
	}

	return project, warnings, nil
}

func classify(path, root string) string {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return ""
	}
	rel = filepath.ToSlash(rel)

	switch {
	case rel == "routes/web.php" || rel == "routes/api.php":
		return "routes"
	case strings.HasPrefix(rel, "app/Http/Controllers/"):
		return "controllers"
	case strings.HasPrefix(rel, "app/Models/"):
		return "models"
	case strings.HasPrefix(rel, "app/Services/"):
		return "services"
	case strings.HasPrefix(rel, "app/Http/Middleware/"):
		return "middleware"
	case strings.HasPrefix(rel, "app/Jobs/"):
		return "jobs"
	case strings.HasPrefix(rel, "app/Events/"):
		return "events"
	case strings.HasPrefix(rel, "app/Listeners/"):
		return "listeners"
	case strings.HasPrefix(rel, "app/Policies/"):
		return "policies"
	default:
		return ""
	}
}
