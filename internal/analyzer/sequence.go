package analyzer

import (
	"fmt"
	"sort"
	"strings"

	"github.com/zam3858/laravel-docgen/internal/model"
)

func buildSequence(project model.Project) []model.SequenceDiagram {
	out := make([]model.SequenceDiagram, 0, len(project.Routes))

	for _, route := range project.Routes {
		participants := []model.Participant{{Alias: "Client", Label: "Client"}}
		if route.Controller != "" {
			participants = append(participants, model.Participant{Alias: route.Controller, Label: route.Controller})
		}

		messages := []model.Message{}
		if route.Controller != "" {
			messages = append(messages,
				model.Message{From: "Client", To: route.Controller, Label: fmt.Sprintf("%s %s", route.Method, route.Path), Type: "sync"},
				model.Message{From: route.Controller, To: "Client", Label: "response", Type: "return"},
			)
		}

		out = append(out, model.SequenceDiagram{
			Type:         "sequence",
			Title:        fmt.Sprintf("%s %s", route.Method, route.Path),
			Participants: participants,
			Messages:     messages,
		})
	}

	sort.Slice(out, func(i, j int) bool {
		return strings.ToLower(out[i].Title) < strings.ToLower(out[j].Title)
	})
	return out
}
