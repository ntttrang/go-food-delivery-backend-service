package restaurantmodel

import "github.com/google/uuid"

type RestaurantDetailReq struct {
	Id uuid.UUID `json:"id"`
}
