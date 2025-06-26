package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharemodel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

// Define DTOs & validate
type UserAddrListReq struct {
	UserAddrSearchDto
	sharemodel.PagingDto
	sharemodel.SortingDto
}

type UserAddrSearchDto struct {
	UserId string `json:"userId" form:"userId"`
	Status string
}

type UserAddrListRes struct {
	Items      []UserAddrSearchResDto `json:"items"`
	Pagination sharemodel.PagingDto   `json:"pagination"`
}

type UserAddrSearchResDto struct {
	Id       uuid.UUID `json:"id"`
	UserId   uuid.UUID `json:"userId"`
	CityId   int       `json:"cityId"`
	CityName string    `json:"cityName"`
	Addr     string    `json:"addr"`
	Lat      *float64  `json:"lat"`
	Lng      *float64  `json:"lng"`
}

// Initilize service
type IListUserAddrRepo interface {
	ListUserAddresses(ctx context.Context, userId uuid.UUID) ([]UserAddrSearchResDto, error)
}

type ListAddrQueryHandler struct {
	userAddrRepo IListUserAddrRepo
}

func NewListAddrQueryHandler(userAddrRepo IListUserAddrRepo) *ListAddrQueryHandler {
	return &ListAddrQueryHandler{
		userAddrRepo: userAddrRepo,
	}
}

// Implement
func (hdl *ListAddrQueryHandler) Execute(ctx context.Context, req UserAddrListReq) (UserAddrListRes, error) {
	userAddrs, err := hdl.userAddrRepo.ListUserAddresses(ctx, uuid.MustParse(req.UserId))

	if err != nil {
		return UserAddrListRes{}, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	var resp UserAddrListRes
	resp.Items = userAddrs
	resp.Pagination = sharemodel.PagingDto{
		Page:  req.Page,
		Limit: req.Limit,
		//Total: total,
	}
	return resp, nil
}
