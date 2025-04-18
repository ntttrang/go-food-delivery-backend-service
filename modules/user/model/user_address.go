package usermodel

import (
	"github.com/google/uuid"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type AddressStatus string

const (
	AddressStatusActive   AddressStatus = "ACTIVE"
	AddressStatusInactive AddressStatus = "INACTIVE"
)

type UserAddress struct {
	Id     uuid.UUID     `json:"id"`
	UserId uuid.UUID     `json:"userId"`
	CityId int           `json:"cityId"`
	Title  *string       `json:"title"`
	Icon   interface{}   `json:"icon"`
	Addr   string        `json:"addr"`
	Lat    *float64      `json:"lat"`
	Lng    *float64      `json:"lng"`
	Status AddressStatus `json:"status"`
	sharedModel.DateDto
}

func (r UserAddress) TableName() string {
	return "user_addresses"
}
