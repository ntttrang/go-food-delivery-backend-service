package restaurantmodel

import (
	"github.com/google/uuid"
)

type RestaurantUpdateReq struct {
	Id  uuid.UUID
	Dto RestaurantUpdateDto
}
