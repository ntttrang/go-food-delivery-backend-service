package restaurantmodel

import (
	"github.com/google/uuid"
)

type MenuItemCreateReq struct {
	RestaurantId uuid.UUID `json:"restaurantId" form:"restaurantId"`
	FoodId       uuid.UUID `json:"foodId" form:"foodId"`
	CategoryId   uuid.UUID `json:"categoryId" form:"categoryId"`
}

func (r MenuItemCreateReq) ConvertToRestaurantFood() *RestaurantFood {
	return &RestaurantFood{
		RestaurantId: r.RestaurantId,
		FoodId:       r.FoodId,
	}
}

func (r MenuItemCreateReq) Validate() error {
	if r.RestaurantId.String() == "" {
		return ErrRestaurantIdRequired
	}
	if r.FoodId.String() == "" {
		return ErrFoodIdRequired
	}
	if r.CategoryId.String() == "" {
		return ErrCategoryIdRequired
	}
	return nil
}
