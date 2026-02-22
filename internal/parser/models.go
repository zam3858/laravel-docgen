package parser

import (
	"fmt"
	"regexp"

	"github.com/zam3858/laravel-docgen/internal/model"
)

var relationRe = regexp.MustCompile(`\$this->(hasOne|hasMany|belongsTo|belongsToMany|morphTo|morphMany|morphOne)\s*\(\s*([A-Za-z0-9_\\]+)::class`)

func parseModels(paths []string) ([]model.Model, []string, error) {
	out := make([]model.Model, 0, len(paths))
	warnings := make([]string, 0)

	for _, path := range paths {
		content, err := readFile(path)
		if err != nil {
			warnings = append(warnings, fmt.Sprintf("read model %s: %v", path, err))
			continue
		}

		rels := make([]model.Relationship, 0)
		for _, m := range relationRe.FindAllStringSubmatch(content, -1) {
			rels = append(rels, model.Relationship{Type: m[1], Related: trimNamespace(m[2])})
		}

		out = append(out, model.Model{
			Name:          parseClass(content, path),
			Namespace:     parseNamespace(content),
			Relationships: rels,
		})
	}

	return out, warnings, nil
}
