package usermodel

import (
	"strings"

	"github.com/google/uuid"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type CreateUserReq struct {
	Email     string      `json:"email"`
	Password  string      `json:"password"`
	FirstName string      `json:"firstName"`
	LastName  string      `json:"lastName"`
	Role      *UserRole   `json:"role"`
	Type      *UserType   `json:"userType"`
	Status    *UserStatus `json:"status"`
	Phone     *string     `json:"phone"`

	Id uuid.UUID `json:"-"`
}

func (r *CreateUserReq) Validate() error {
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

func (r CreateUserReq) ConvertToUser() *User {
	return &User{
		Email:     r.Email,
		Password:  r.Password,
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Role:      *r.Role,
		Type:      *r.Type,
		Status:    *r.Status,
		Phone:     *r.Phone,
	}
}
