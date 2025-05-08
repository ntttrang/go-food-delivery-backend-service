package datatype

import "github.com/google/uuid"

type Requester interface {
	Subject() uuid.UUID
	GetRole() uuid.UUID
}
type requester struct {
	userId uuid.UUID
	role   string
}

func NewRequester(userId string) *requester {
	return &requester{
		userId: uuid.MustParse(userId),
	}
}

func (r *requester) Subject() uuid.UUID {
	return r.userId
}

func (r *requester) GetRole() string {
	return r.role
}
