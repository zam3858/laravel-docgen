package parser

import (
	"fmt"
	"regexp"

	"github.com/zam3858/laravel-docgen/internal/model"
)

var handleRe = regexp.MustCompile(`public\s+function\s+handle\s*\(\s*([A-Za-z0-9_\\]+)\s+\$`)

func parseEvents(paths []string) ([]model.Event, []string, error) {
	out := make([]model.Event, 0, len(paths))
	warnings := make([]string, 0)

	for _, path := range paths {
		content, err := readFile(path)
		if err != nil {
			warnings = append(warnings, fmt.Sprintf("read event %s: %v", path, err))
			continue
		}
		out = append(out, model.Event{Name: parseClass(content, path), Namespace: parseNamespace(content)})
	}

	return out, warnings, nil
}

func parseListeners(paths []string) ([]model.Listener, []string, error) {
	out := make([]model.Listener, 0, len(paths))
	warnings := make([]string, 0)

	for _, path := range paths {
		content, err := readFile(path)
		if err != nil {
			warnings = append(warnings, fmt.Sprintf("read listener %s: %v", path, err))
			continue
		}
		handles := make([]string, 0)
		for _, m := range handleRe.FindAllStringSubmatch(content, -1) {
			handles = append(handles, trimNamespace(m[1]))
		}
		out = append(out, model.Listener{Name: parseClass(content, path), Namespace: parseNamespace(content), Handles: handles})
	}

	return out, warnings, nil
}
