package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	usermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/user/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type IupdateRepo interface {
	FindById(ctx context.Context, id uuid.UUID) (*usermodel.User, error)
	Update(ctx context.Context, req *usermodel.UpdateUserReq) error
}

type UpdateCommandHandler struct {
	userRepo IupdateRepo
}

func NewUpdateCommandHandler(userRepo IupdateRepo) *UpdateCommandHandler {
	return &UpdateCommandHandler{userRepo: userRepo}
}

func (hdl *UpdateCommandHandler) Execute(ctx context.Context, req usermodel.UpdateUserReq) error {
	if err := req.Validate(); err != nil {
		return datatype.ErrBadRequest.WithWrap(err).WithDebug(err.Error())
	}

	existUser, err := hdl.userRepo.FindById(ctx, req.Id)

	if err != nil {
		if errors.Is(err, usermodel.ErrUserNotFound) {
			return datatype.ErrNotFound.WithDebug(usermodel.ErrUserNotFound.Error())
		}

		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	if existUser.Status == sharedModel.StatusDelete {
		return datatype.ErrNotFound.WithError(usermodel.ErrUserDeletedOrBanned.Error())
	}

	if err := hdl.userRepo.Update(ctx, &req); err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	return nil
}
