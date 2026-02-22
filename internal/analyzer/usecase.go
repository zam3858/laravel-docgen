package analyzer

import (
	"fmt"
	"sort"
	"strings"

	"github.com/zam3858/laravel-docgen/internal/model"
)

func buildUseCases(project model.Project) []model.UseCaseDiagram {
	if len(project.Routes) == 0 {
		return nil
	}

	usecases := make([]model.UseCase, 0, len(project.Routes))
	rels := make([]model.UseCaseRelationship, 0, len(project.Routes))

	for i, r := range project.Routes {
		id := fmt.Sprintf("UC%d", i+1)
		label := fmt.Sprintf("%s %s", r.Method, r.Path)
		usecases = append(usecases, model.UseCase{ID: id, Label: label})
		rels = append(rels, model.UseCaseRelationship{Actor: "Client", UseCase: id, Type: "association"})
	}

	sort.Slice(usecases, func(i, j int) bool { return strings.ToLower(usecases[i].Label) < strings.ToLower(usecases[j].Label) })
	sort.Slice(rels, func(i, j int) bool { return rels[i].UseCase < rels[j].UseCase })

	return []model.UseCaseDiagram{{
		Type:          "usecase",
		Title:         "Application Use Cases",
		Actors:        []model.Actor{{Name: "Client"}},
		UseCases:      usecases,
		Relationships: rels,
	}}
}
