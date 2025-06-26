package usermodel

import (
	"github.com/google/uuid"
	sharedatatype "github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharedmodel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type UserAddress struct {
	Id     uuid.UUID                  `json:"id"`
	UserId uuid.UUID                  `json:"userId"`
	CityId int                        `json:"cityId"`
	Title  *string                    `json:"title"`
	Icon   string                     `json:"icon"`
	Addr   string                     `json:"addr"`
	Lat    *float64                   `json:"lat"`
	Lng    *float64                   `json:"lng"`
	Status sharedatatype.RecordStatus `json:"status"`
	sharedmodel.DateDto
}

func (r UserAddress) TableName() string {
	return "user_addresses"
}
