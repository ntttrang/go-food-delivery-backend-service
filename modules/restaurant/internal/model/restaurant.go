package restaurantmodel

import (
	"encoding/json"

	"github.com/google/uuid"
	sharedmodel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type Restaurant struct {
	Id               uuid.UUID       `gorm:"column:id;"`
	OwnerId          string          `gorm:"column:owner_id;"`
	Name             string          `gorm:"column:name;"`
	Addr             string          `gorm:"column:addr;"`
	CityId           string          `gorm:"column:city_id;"`
	Lat              float64         `gorm:"column:lat;"`
	Lng              float64         `gorm:"column:lng;"`
	Cover            json.RawMessage `gorm:"column:cover;"` // json
	Logo             json.RawMessage `gorm:"column:logo;"`  // json
	ShippingFeePerKm float64         `gorm:"column:shipping_fee_per_km;"`
	Status           string          `gorm:"column:status;"`
	sharedmodel.AbstractInfo
}
