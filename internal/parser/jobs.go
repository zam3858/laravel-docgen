package parser

import (
	"fmt"

	"github.com/zam3858/laravel-docgen/internal/model"
)

func parseJobs(paths []string) ([]model.Job, []string, error) {
	out := make([]model.Job, 0, len(paths))
	warnings := make([]string, 0)

	for _, path := range paths {
		content, err := readFile(path)
		if err != nil {
			warnings = append(warnings, fmt.Sprintf("read job %s: %v", path, err))
			continue
		}
		out = append(out, model.Job{Name: parseClass(content, path), Namespace: parseNamespace(content)})
	}

	return out, warnings, nil
}
