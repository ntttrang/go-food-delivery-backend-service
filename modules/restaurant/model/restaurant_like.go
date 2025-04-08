package restaurantmodel

import (
	"github.com/google/uuid"
	sharedmodel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type RestaurantLike struct {
	RestaurantID uuid.UUID `gorm:"column:restaurant_id" json:"restaurantId" form:"restaurantId"`
	UserID       uuid.UUID `gorm:"column:user_id" json:"userId" form:"userId"`
	sharedmodel.DateDto
}

func (r RestaurantLike) Validate() error {
	if r.RestaurantID.String() == "" {
		return ErrRestaurantIdRequired
	}

	if r.UserID.String() == "" {
		return ErrUserIdRequired
	}

	return nil
}
