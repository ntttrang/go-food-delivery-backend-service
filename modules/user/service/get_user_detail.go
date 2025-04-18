package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	usermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/user/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type IGetDetailQueryRepo interface {
	FindById(ctx context.Context, id uuid.UUID) (*usermodel.User, error)
}

type GetDetailQueryHandler struct {
	repo IGetDetailQueryRepo
}

func NewGetDetailQueryHandler(repo IGetDetailQueryRepo) *GetDetailQueryHandler {
	return &GetDetailQueryHandler{repo: repo}
}

func (hdl *GetDetailQueryHandler) Execute(ctx context.Context, req usermodel.UserDetailReq) (usermodel.UserSearchResDto, error) {
	user, err := hdl.repo.FindById(ctx, req.Id)

	if err != nil {
		if errors.Is(err, usermodel.ErrUserNotFound) {
			return usermodel.UserSearchResDto{}, datatype.ErrNotFound.WithDebug(usermodel.ErrUserNotFound.Error())
		}
		return usermodel.UserSearchResDto{}, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	if user.Status == sharedModel.StatusDelete {
		return usermodel.UserSearchResDto{}, datatype.ErrDeleted.WithError(usermodel.ErrUserDeletedOrBanned.Error())
	}

	var resp usermodel.UserSearchResDto
	if err := copier.Copy(&resp, &user); err != nil {
		return usermodel.UserSearchResDto{}, datatype.ErrInternalServerError.WithWrap(errors.New("copier libraries failed"))
	}
	return resp, nil
}
