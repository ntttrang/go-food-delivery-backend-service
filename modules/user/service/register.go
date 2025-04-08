package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	usermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/user/model"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
	"golang.org/x/crypto/bcrypt"
)

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

func (hdl *RegisterUserCommandHandler) Execute(ctx context.Context, req usermodel.RegisterUserReq) error {
	if err := req.Validate(); err != nil {
		return err
	}

	existUser, err := hdl.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return err
	}

	if existUser != nil {
		if existUser.Status == usermodel.StatusBanned || existUser.Status == usermodel.StatusDeleted {
			return usermodel.ErrUserDeletedOrBanned
		}

	}

	salt, _ := sharedModel.RandomStr(16)
	saltPass := fmt.Sprintf("%s.%s", salt, req.Password)
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(saltPass), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("generate from password failed")
	}

	id, _ := uuid.NewV7()
	now := time.Now().UTC()

	var user usermodel.User
	if err = copier.Copy(&user, &req); err != nil {
		fmt.Println(err)
		return err
	}

	user.Id = id
	user.Password = string(hashPassword)
	user.Salt = salt
	user.Status = usermodel.StatusActive
	user.Type = usermodel.TypeEmailPassword
	user.Role = usermodel.RoleAdmin
	user.CreatedAt = &now
	user.UpdatedAt = &now

	if err := hdl.repo.Insert(ctx, &user); err != nil {
		return err
	}

	req.Id = id
	return nil

}
