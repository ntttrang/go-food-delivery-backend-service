package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	usermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/user/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharemodel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
	"golang.org/x/crypto/bcrypt"
)

// Define DTOs & validate
type CreateUserReq struct {
	Email     string               `json:"email"`
	Password  string               `json:"password"`
	FirstName string               `json:"firstName"`
	LastName  string               `json:"lastName"`
	Role      *datatype.UserRole   `json:"role"`
	Type      *datatype.UserType   `json:"userType"`
	Status    *datatype.UserStatus `json:"status"`
	Phone     *string              `json:"phone"`

	Id        uuid.UUID          `json:"-"`
	Requester datatype.Requester `json:"-"`
}

func (r *CreateUserReq) Validate() error {
	r.Email = strings.TrimSpace(r.Email)
	r.Password = strings.TrimSpace(r.Password)
	r.FirstName = strings.TrimSpace(r.FirstName)
	r.LastName = strings.TrimSpace(r.LastName)

	if r.Email == "" {
		return usermodel.ErrEmailRequired
	}

	if !sharemodel.ValidateEmail(r.Email) {
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

	// Authorization check
	if req.Requester == nil {
		return datatype.ErrUnauthorized.WithDebug("requester information required")
	}

	// Check permissions:
	// 1. Users can update their own profile
	// 2. Admins can update any user's profile
	isOwnProfile := req.Requester.Subject() == req.Id

	// Check if requester is admin
	isAdmin := false
	if user, ok := req.Requester.(*IntrospectRes); ok {
		isAdmin = user.Role == datatype.RoleAdmin
	}

	if !isOwnProfile && !isAdmin {
		return datatype.ErrForbidden.WithDebug(usermodel.ErrPermission.Error())
	}

	user := req.ConvertToUser()
	user.Id, _ = uuid.NewV7()
	salt, _ := sharemodel.RandomStr(16)
	saltPass := fmt.Sprintf("%s.%s", salt, req.Password)
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(saltPass), bcrypt.DefaultCost)
	if err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}
	user.Password = string(hashPassword)
	user.Salt = salt
	user.Status = datatype.StatusActive // Always set Active Status when insert
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
