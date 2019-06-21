package bellbox

import "time"

type User struct {
	User string
	Password string
	Admin bool `json:-` // ignored when submitted by clients
}

type UserToken struct {
	User string
	Token string
	Timestamp time.Time `gorm:"type:time"`
}

type Bell struct {
	Id string
	User string `json:-` // ignore clients
	Name string
	Type string
	Key string
	Enabled bool `json:-` // ignore clients
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
	Urgent bool
}

type FilterRequest struct {
	Type string
}
