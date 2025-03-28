package restaurantmodel

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type RestaurantDetailRes struct {
	Id               uuid.UUID       `json:"id"`
	Name             string          `json:"name"`
	Addr             string          `json:"addr"`
	Logo             json.RawMessage `json:"logo"`
	ShippingFeePerKm float64         `json:"shippingFeePerKm"`
	Status           string          `json:"status"`
	//CategoryName     string          `json:"categoryName"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
}
