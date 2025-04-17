package usermodel

import (
	"strings"

	"github.com/google/uuid"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type UserAddrListReq struct {
	UserAddrSearchDto
	sharedModel.PagingDto
	sharedModel.SortingDto
}

type UserAddrSearchDto struct {
	UserId string `json:"userId" form:"userId"`
	Status string
}

type UserAddrDetailReq struct {
	Id uuid.UUID `json:"-"`
}

type CreateUserAddrReq struct {
	CityId int         `json:"cityId"`
	Title  *string     `json:"title"`
	Icon   interface{} `json:"icon"`
	Addr   string      `json:"addr"`
	Lat    *float64    `json:"lat"`
	Lng    *float64    `json:"lng"`

	UserId uuid.UUID `json:"-"`
	Id     uuid.UUID `json:"-"`
}

func (r CreateUserAddrReq) TableName() string {
	return UserAddress{}.TableName()
}

func (r *CreateUserAddrReq) Validate() error {
	r.Addr = strings.TrimSpace(r.Addr)
	if r.Addr == "" {
		return ErrAddrRequired
	}
	return nil
}

type UpdateUserAddrReq struct {
	CityId int         `json:"cityId"`
	Title  *string     `json:"title"`
	Icon   interface{} `json:"icon"`
	Addr   string      `json:"addr"`
	Lat    *float64    `json:"lat"`
	Lng    *float64    `json:"lng"`
	Status string      `json:"status"`

	Id uuid.UUID `json:"-"`
}
