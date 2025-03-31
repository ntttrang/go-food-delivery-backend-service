package restaurantmodel

import (
	"encoding/json"

	"github.com/google/uuid"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type RestaurantSearchRes struct {
	Items      []RestaurantSearchResDto `json:"items"`
	Pagination sharedModel.PagingDto    `json:"pagination"`
}

type RestaurantSearchResDto struct {
	Id               uuid.UUID       `json:"id"`
	Name             string          `json:"name"`
	Addr             string          `json:"addr"`
	Logo             json.RawMessage `json:"logo"`
	ShippingFeePerKm float64         `json:"shippingFeePerKm"`
	Status           string          `json:"status"`
	//CategoryName     string          `json:"categoryName"`
}
