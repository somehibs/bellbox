package bellbox

import "time"

type Message struct {
	Target string
	Title string
	Message string
	Url string
	Timestamp time.Time
	Priority string
	Sender string
}

type User struct {
	User string
	Password string
	Admin bool `json:"-"` // ignored when submitted by clients
}

type UserToken struct {
	User string
	Token string
	Timestamp time.Time
}

type Bell struct {
	Id string // ignore clients but allow it to be returned
	User string `json:"-"` // ignore clients
	Name string
	Type string
	Key string
	Enabled bool `json:"-"` // ignore clients
}

type UserReply struct {
	Token string
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
	Token string `json:"-"`
	Urgent bool
	RequestState int
}

type FilterRequest struct {
	Type string
}
