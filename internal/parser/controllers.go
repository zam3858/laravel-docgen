package parser

import (
	"fmt"

	"github.com/zam3858/laravel-docgen/internal/model"
)

func parseControllers(paths []string) ([]model.Controller, []string, error) {
	out := make([]model.Controller, 0, len(paths))
	warnings := make([]string, 0)

	for _, path := range paths {
		content, err := readFile(path)
		if err != nil {
			warnings = append(warnings, fmt.Sprintf("read controller %s: %v", path, err))
			continue
		}

		actions := make([]model.Action, 0)
		for _, m := range methodRe.FindAllStringSubmatch(content, -1) {
			actions = append(actions, model.Action{Name: m[1]})
		}

		out = append(out, model.Controller{
			Name:      parseClass(content, path),
			Namespace: parseNamespace(content),
			Actions:   actions,
		})
	}

	return out, warnings, nil
}
