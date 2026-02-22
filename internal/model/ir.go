package model

type Project struct {
	Routes      []Route
	Controllers []Controller
	Models      []Model
	Services    []Service
	Middleware  []Middleware
	Jobs        []Job
	Events      []Event
	Listeners   []Listener
	Policies    []Policy
}

type Route struct {
	Method     string
	Path       string
	Middleware []string
	Controller string
	Action     string
}

type Controller struct {
	Name      string
	Namespace string
	Actions   []Action
}

type Action struct {
	Name         string
	Dependencies []string
}

type Model struct {
	Name          string
	Namespace     string
	Fields        []Field
	Relationships []Relationship
}

type Field struct {
	Name string
	Type string
}

type Relationship struct {
	Type    string
	Related string
	Name    string
}

type Service struct {
	Name      string
	Namespace string
}

type Middleware struct {
	Name      string
	Namespace string
}

type Job struct {
	Name      string
	Namespace string
}

type Event struct {
	Name      string
	Namespace string
}

type Listener struct {
	Name      string
	Namespace string
	Handles   []string
}

type Policy struct {
	Name      string
	Namespace string
}
