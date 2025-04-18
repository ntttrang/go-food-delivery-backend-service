package restaurantmodel

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	sharedmodel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
	"gorm.io/gorm"
)

type Restaurant struct {
	Id               uuid.UUID       `gorm:"column:id;"`
	OwnerId          uuid.UUID       `gorm:"column:owner_id;"`
	Name             string          `gorm:"column:name;"`
	Addr             string          `gorm:"column:addr;"`
	CityId           int             `gorm:"column:city_id;"`
	Lat              float64         `gorm:"column:lat;"`
	Lng              float64         `gorm:"column:lng;"`
	Cover            json.RawMessage `gorm:"column:cover;"` // json
	Logo             json.RawMessage `gorm:"column:logo;"`  // json
	ShippingFeePerKm float64         `gorm:"column:shipping_fee_per_km;"`
	Status           string          `gorm:"column:status;"`
	sharedmodel.DateDto
}

func (Restaurant) TableName() string {
	return "restaurants"
}

func (r *Restaurant) BeforeCreate(tx *gorm.DB) {
	//currentTime := time.Now().UTC()
	currentTime := time.Date(2021, time.March, 28, 12, 31, 0, 0, time.UTC)
	r.CreatedAt = &currentTime
	r.UpdatedAt = &currentTime
}
