package analyzer

import (
	"fmt"
	"sort"
	"strings"

	"github.com/zam3858/laravel-docgen/internal/model"
)

func buildObjects(project model.Project) []model.ObjectDiagram {
	out := make([]model.ObjectDiagram, 0, len(project.Models))
	for _, m := range project.Models {
		rels := make([]model.ObjectRelationship, 0, len(m.Relationships))
		for _, r := range m.Relationships {
			rels = append(rels, model.ObjectRelationship{From: m.Name, To: r.Related, Label: r.Type})
		}
		sort.Slice(rels, func(i, j int) bool {
			li := rels[i].From + rels[i].To + rels[i].Label
			lj := rels[j].From + rels[j].To + rels[j].Label
			return strings.ToLower(li) < strings.ToLower(lj)
		})

		out = append(out, model.ObjectDiagram{
			Type:          "object",
			Title:         fmt.Sprintf("%s Model", m.Name),
			Objects:       []model.ObjectEntry{{Name: m.Name, Fields: m.Fields}},
			Relationships: rels,
		})
	}

	sort.Slice(out, func(i, j int) bool {
		return strings.ToLower(out[i].Title) < strings.ToLower(out[j].Title)
	})
	return out
}
