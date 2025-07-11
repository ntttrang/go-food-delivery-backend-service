package service

import (
	"context"
	"errors"
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

func (r RegisterUserReq) ConvertToUser() *usermodel.User {
	return &usermodel.User{
		Email:     r.Email,
		Password:  r.Password,
		FirstName: r.FirstName,
		LastName:  r.LastName,
	}
}

// Initilize service
type IRegisterRepo interface {
	FindByEmail(ctx context.Context, email string) (*usermodel.User, error)
	Insert(ctx context.Context, user *usermodel.User) error
}

type RegisterUserCommandHandler struct {
	repo IRegisterRepo
}

func NewRegisterUserCommandHandler(repo IRegisterRepo) *RegisterUserCommandHandler {
	return &RegisterUserCommandHandler{repo: repo}
}

// Implement
func (hdl *RegisterUserCommandHandler) Execute(ctx context.Context, req *RegisterUserReq) error {
	if err := req.Validate(); err != nil {
		return datatype.ErrBadRequest.WithWrap(err).WithDebug(err.Error())
	}

	existUser, err := hdl.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		if !errors.Is(err, usermodel.ErrUserNotFound) {
			return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
		}

	}

	if existUser != nil {
		if existUser.Status == datatype.StatusBanned || existUser.Status == datatype.StatusDeleted {
			return usermodel.ErrUserDeletedOrBanned
		}

	}

	salt, _ := sharemodel.RandomStr(16)
	saltPass := fmt.Sprintf("%s.%s", salt, req.Password)
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(saltPass), bcrypt.DefaultCost)
	if err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	id, _ := uuid.NewV7()
	now := time.Now().UTC()

	var user = usermodel.User{
		Id:        id,
		Email:     req.Email,
		Password:  string(hashPassword),
		LastName:  req.LastName,
		FirstName: req.FirstName,
		Salt:      salt,
		Status:    datatype.StatusActive,
		Type:      datatype.TypeEmailPassword,
		Role:      datatype.RoleUser,
		DateDto: sharemodel.DateDto{
			CreatedAt: &now,
			UpdatedAt: &now,
		},
	}

	if err := hdl.repo.Insert(ctx, &user); err != nil {
		return err
	}

	req.Id = id
	return nil

}
