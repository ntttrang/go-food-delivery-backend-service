package datatype

type Requester interface {
	Subject() string
	Role() string
}
type requester struct {
	userId string
	role   string
}

func NewRequester(userId string) *requester {
	return &requester{
		userId: userId,
	}
}

func (r *requester) Subject() string {
	return r.userId
}

func (r *requester) Role() string {
	return r.role
}
