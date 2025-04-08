package usermodel

import (
	"strings"

	"github.com/google/uuid"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type RegisterUserReq struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`

	Id uuid.UUID `json:"-"`
}

func (r *RegisterUserReq) Validate() error {
	r.Email = strings.TrimSpace(r.Email)
	r.Password = strings.TrimSpace(r.Password)
	r.FirstName = strings.TrimSpace(r.FirstName)
	r.LastName = strings.TrimSpace(r.LastName)

	if r.Email == "" {
		return ErrEmailRequired
	}

	if !sharedModel.ValidateEmail(r.Email) {
		return ErrEmailInvalid
	}

	if len(r.Password) <= 6 {
		return ErrPasswordInvalid
	}

	if r.FirstName == "" {
		return ErrFirstNameRequired
	}

	if r.LastName == "" {
		return ErrLastNameRequired
	}

	return nil
}
