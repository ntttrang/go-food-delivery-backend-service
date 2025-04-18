package service

import (
	"context"

	"github.com/google/uuid"
	usermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/user/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type IListUserAddrRepo interface {
	ListUserAddresses(ctx context.Context, userId uuid.UUID) ([]usermodel.UserAddrSearchResDto, error)
}

type ListAddrQueryHandler struct {
	userAddrRepo IListUserAddrRepo
}

func NewListAddrQueryHandler(userAddrRepo IListUserAddrRepo) *ListAddrQueryHandler {
	return &ListAddrQueryHandler{
		userAddrRepo: userAddrRepo,
	}
}

func (hdl *ListAddrQueryHandler) Execute(ctx context.Context, req usermodel.UserAddrListReq) (usermodel.UserAddrListRes, error) {
	userAddrs, err := hdl.userAddrRepo.ListUserAddresses(ctx, uuid.MustParse(req.UserId))

	if err != nil {
		return usermodel.UserAddrListRes{}, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	var resp usermodel.UserAddrListRes
	resp.Items = userAddrs
	resp.Pagination = sharedModel.PagingDto{
		Page:  req.Page,
		Limit: req.Limit,
		//Total: total,
	}
	return resp, nil
}
