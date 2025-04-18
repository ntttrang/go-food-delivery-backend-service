package usermodel

import "github.com/google/uuid"

type UserDetailReq struct {
	Id uuid.UUID `json:"id"`
}
