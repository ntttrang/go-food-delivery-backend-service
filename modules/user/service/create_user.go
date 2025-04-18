package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	usermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/user/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type ICreateRepo interface {
	Insert(ctx context.Context, data *usermodel.User) error
}

type CreateCommandHandler struct {
	userRepo ICreateRepo
}

func NewCreateCommandHandler(userRepo ICreateRepo) *CreateCommandHandler {
	return &CreateCommandHandler{userRepo: userRepo}
}

func (s *CreateCommandHandler) Execute(ctx context.Context, req *usermodel.CreateUserReq) error {
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
