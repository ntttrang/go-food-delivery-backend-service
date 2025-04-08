package usermodel

import (
	"strings"

	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type AuthenticateReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *AuthenticateReq) Validate() error {
	r.Email = strings.TrimSpace(r.Email)
	r.Password = strings.TrimSpace(r.Password)

	if r.Email == "" {
		return ErrEmailRequired
	}

	if r.Password == "" {
		return ErrPasswordInvalid
	}

	if !sharedModel.ValidateEmail(r.Email) {
		return ErrEmailInvalid
	}

	if len(r.Password) <= 6 {
		return ErrPasswordInvalid
	}

	return nil
}
