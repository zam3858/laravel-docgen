package parser

import (
	"fmt"

	"github.com/zam3858/laravel-docgen/internal/model"
)

func parsePolicies(paths []string) ([]model.Policy, []string, error) {
	out := make([]model.Policy, 0, len(paths))
	warnings := make([]string, 0)

	for _, path := range paths {
		content, err := readFile(path)
		if err != nil {
			warnings = append(warnings, fmt.Sprintf("read policy %s: %v", path, err))
			continue
		}
		out = append(out, model.Policy{Name: parseClass(content, path), Namespace: parseNamespace(content)})
	}

	return out, warnings, nil
}
