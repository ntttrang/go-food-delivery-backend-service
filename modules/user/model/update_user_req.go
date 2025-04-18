package usermodel

import (
	"strings"

	"github.com/google/uuid"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type UpdateUserReq struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Phone     string `json:"phone"`
	Avatar    string `json:"avatar"`
	Status    string `json:"status"`

	Id uuid.UUID `json:"_"`
}

func (UpdateUserReq) TableName() string {
	return User{}.TableName()
}

func (r *UpdateUserReq) Validate() error {
	r.FirstName = strings.TrimSpace(r.FirstName)
	r.LastName = strings.TrimSpace(r.LastName)
	r.Phone = strings.TrimSpace(r.Phone)

	if r.Id.String() == "" {
		return ErrIdRequired
	}

	if r.Phone != "" && !sharedModel.ValidatePhoneNumber(r.Phone) {
		return ErrInvalidPhoneNumber
	}

	return nil
}
