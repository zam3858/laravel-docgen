package parser

import (
	"fmt"

	"github.com/zam3858/laravel-docgen/internal/model"
)

func parseMiddleware(paths []string) ([]model.Middleware, []string, error) {
	out := make([]model.Middleware, 0, len(paths))
	warnings := make([]string, 0)

	for _, path := range paths {
		content, err := readFile(path)
		if err != nil {
			warnings = append(warnings, fmt.Sprintf("read middleware %s: %v", path, err))
			continue
		}
		out = append(out, model.Middleware{Name: parseClass(content, path), Namespace: parseNamespace(content)})
	}

	return out, warnings, nil
}
