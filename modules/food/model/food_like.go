package foodmodel

import (
	"github.com/google/uuid"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type FoodLike struct {
	UserId uuid.UUID `gorm:"column:user_id"`
	FoodId uuid.UUID `gorm:"column:food_id"`
	sharedModel.DateDto
}

func (FoodLike) TableName() string {
	return "food_likes"
}

func (r FoodLike) Validate() error {
	if r.FoodId.String() == "" {
		return ErrFoodIdRequired
	}

	if r.UserId.String() == "" {
		return ErrUserIdRequired
	}

	return nil
}
