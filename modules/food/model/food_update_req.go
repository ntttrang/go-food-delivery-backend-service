package foodmodel

import (
	"github.com/google/uuid"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type FoodUpdateReq struct {
	// Use pointer to accept empty string
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Status       *string `json:"status"`
	RestaurantId *string `json:"restaurantId"` // Can be empty or missing if data type = string. Otherwise, uuid.UUID isn't
	CategoryId   *string `json:"categoryId"`   // Can be empty or missing if data type = string. Otherwise, uuid.UUID isn't

	Id uuid.UUID `json:"-"`
}

func (FoodUpdateReq) TableName() string {
	return Food{}.TableName()
}

func (c FoodUpdateReq) Validate() error {
	if c.Status != nil && *c.Status != sharedModel.StatusActive && *c.Status != sharedModel.StatusDelete && *c.Status != sharedModel.StatusInactive {
		return ErrFoodStatusInvalid
	}
	return nil
}
