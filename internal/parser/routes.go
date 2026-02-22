package parser

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/zam3858/laravel-docgen/internal/model"
)

var routeRe = regexp.MustCompile(`Route::(get|post|put|patch|delete|options)\s*\(\s*['\"]([^'\"]+)['\"]\s*,\s*\[\s*([A-Za-z0-9_\\]+)::class\s*,\s*['\"]([A-Za-z0-9_]+)['\"]\s*\]\s*\)`)

func parseRoutes(paths []string) ([]model.Route, []string, error) {
	out := make([]model.Route, 0)
	warnings := make([]string, 0)

	for _, path := range paths {
		content, err := readFile(path)
		if err != nil {
			warnings = append(warnings, fmt.Sprintf("read route file %s: %v", path, err))
			continue
		}
		matches := routeRe.FindAllStringSubmatch(content, -1)
		for _, m := range matches {
			out = append(out, model.Route{
				Method:     strings.ToUpper(m[1]),
				Path:       m[2],
				Middleware: nil,
				Controller: trimNamespace(m[3]),
				Action:     m[4],
			})
		}
	}

	return out, warnings, nil
}

func trimNamespace(s string) string {
	if idx := strings.LastIndex(s, `\\`); idx >= 0 {
		return s[idx+1:]
	}
	return s
}
