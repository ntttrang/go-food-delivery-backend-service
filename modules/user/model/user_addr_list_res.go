package usermodel

import (
	"github.com/google/uuid"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type UserAddrListRes struct {
	Items      []UserAddrSearchResDto `json:"items"`
	Pagination sharedModel.PagingDto  `json:"pagination"`
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
