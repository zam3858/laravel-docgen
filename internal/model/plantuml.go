package model

type Participant struct {
	Alias string `json:"alias"`
	Label string `json:"label"`
}

type Message struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Label string `json:"label"`
	Type  string `json:"type"`
}

type SequenceDiagram struct {
	Type         string        `json:"type"`
	Title        string        `json:"title"`
	Participants []Participant `json:"participants"`
	Messages     []Message     `json:"messages"`
}

type ObjectEntry struct {
	Name   string  `json:"name"`
	Fields []Field `json:"fields"`
}

type ObjectRelationship struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Label string `json:"label"`
}

type ObjectDiagram struct {
	Type          string               `json:"type"`
	Title         string               `json:"title"`
	Objects       []ObjectEntry        `json:"objects"`
	Relationships []ObjectRelationship `json:"relationships"`
}

type Actor struct {
	Name string `json:"name"`
}

type UseCase struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

type UseCaseRelationship struct {
	Actor   string `json:"actor"`
	UseCase string `json:"usecase"`
	Type    string `json:"type"`
}

type UseCaseDiagram struct {
	Type          string                `json:"type"`
	Title         string                `json:"title"`
	Actors        []Actor               `json:"actors"`
	UseCases      []UseCase             `json:"usecases"`
	Relationships []UseCaseRelationship `json:"relationships"`
}

type Analysis struct {
	Sequence []SequenceDiagram
	Objects  []ObjectDiagram
	UseCase  []UseCaseDiagram
}
