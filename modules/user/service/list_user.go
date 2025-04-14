package service

import (
	"context"

	usermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/user/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type IListUserRepo interface {
	FindUsers(ctx context.Context, req usermodel.UserListReq) ([]usermodel.UserSearchResDto, int64, error)
}

type ListQueryHandler struct {
	userRepo IListUserRepo
}

func NewListQueryHandler(userRepo IListUserRepo) *ListQueryHandler {
	return &ListQueryHandler{
		userRepo: userRepo,
	}
}

func (hdl *ListQueryHandler) Execute(ctx context.Context, req usermodel.UserListReq) (usermodel.UserListRes, error) {
	users, total, err := hdl.userRepo.FindUsers(ctx, req)

	if err != nil {
		return usermodel.UserListRes{}, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	var resp usermodel.UserListRes
	resp.Items = users
	resp.Pagination = sharedModel.PagingDto{
		Page:  req.Page,
		Limit: req.Limit,
		Total: total,
	}
	return resp, nil
}
