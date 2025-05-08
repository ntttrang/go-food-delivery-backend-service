package service

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	usermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/user/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

// Define DTOs & validate
type CreateUserReq struct {
	Email     string                `json:"email"`
	Password  string                `json:"password"`
	FirstName string                `json:"firstName"`
	LastName  string                `json:"lastName"`
	Role      *usermodel.UserRole   `json:"role"`
	Type      *usermodel.UserType   `json:"userType"`
	Status    *usermodel.UserStatus `json:"status"`
	Phone     *string               `json:"phone"`

	Id uuid.UUID `json:"-"`
}

func (r *CreateUserReq) Validate() error {
	r.Email = strings.TrimSpace(r.Email)
	r.Password = strings.TrimSpace(r.Password)
	r.FirstName = strings.TrimSpace(r.FirstName)
	r.LastName = strings.TrimSpace(r.LastName)

	if r.Email == "" {
		return usermodel.ErrEmailRequired
	}

	if !sharedModel.ValidateEmail(r.Email) {
		return usermodel.ErrEmailInvalid
	}

	if len(r.Password) <= 6 {
		return usermodel.ErrPasswordInvalid
	}

	if r.FirstName == "" {
		return usermodel.ErrFirstNameRequired
	}

	if r.LastName == "" {
		return usermodel.ErrLastNameRequired
	}

	return nil
}

func (r CreateUserReq) ConvertToUser() *usermodel.User {
	return &usermodel.User{
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

// Initilize service
type ICreateRepo interface {
	Insert(ctx context.Context, data *usermodel.User) error
}

type CreateCommandHandler struct {
	userRepo ICreateRepo
}

func NewCreateCommandHandler(userRepo ICreateRepo) *CreateCommandHandler {
	return &CreateCommandHandler{userRepo: userRepo}
}

// Implement
func (s *CreateCommandHandler) Execute(ctx context.Context, req *CreateUserReq) error {
	if err := req.Validate(); err != nil {
		return datatype.ErrBadRequest.WithWrap(err).WithDebug(err.Error())
	}

	user := req.ConvertToUser()
	user.Id, _ = uuid.NewV7()
	user.Status = sharedModel.StatusActive // Always set Active Status when insert
	now := time.Now().UTC()
	user.CreatedAt = &now
	user.UpdatedAt = &now

	if err := s.userRepo.Insert(ctx, user); err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	// set data to response
	req.Id = user.Id

	return nil
}
