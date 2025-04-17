package foodmodel

import (
	"github.com/google/uuid"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type FoodRatings struct {
	Id      uuid.UUID `gorm:"column:id"`
	UserId  uuid.UUID `gorm:"column:user_id"`
	FoodId  uuid.UUID `gorm:"column:food_id"`
	Point   float64   `gorm:"column:point"`
	Comment string    `gorm:"column:comment"`
	Status  string    `gorm:"column:status"`
	sharedModel.DateDto
}

func (FoodRatings) TableName() string {
	return "food_ratings"
}
