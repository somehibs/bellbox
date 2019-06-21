package bellbox

type User struct {
	User string
	Password string
	Admin bool `json:-` // ignored when submitted by clients
	Token string `json:-` // ignored when submitted by clients
}

type Bell struct {
	Id string
	Name string
	Type string
	Key string
	Enabled bool
}

type Id struct {
	Id string
}

type DeleteBellRequest struct {
	Id string
	Secret string
}

type Bellringer struct {
	Target string
	Name string
	Urgent bool
}

type FilterRequest struct {
	Type string
}
