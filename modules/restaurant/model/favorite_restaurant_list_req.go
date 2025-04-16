package restaurantmodel

import (
	"github.com/google/uuid"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type FavoriteRestaurantListReq struct {
	UserId uuid.UUID `json:"-" form:"-"`
	sharedModel.PagingDto
}
