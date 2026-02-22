package analyzer

import "github.com/zam3858/laravel-docgen/internal/model"

type Analyzer struct{}

func New() *Analyzer {
	return &Analyzer{}
}

func (a *Analyzer) Analyze(project model.Project) model.Analysis {
	return model.Analysis{
		Sequence: buildSequence(project),
		Objects:  buildObjects(project),
		UseCase:  buildUseCases(project),
	}
}
