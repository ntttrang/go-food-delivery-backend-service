package foodmodel

import (
	"github.com/google/uuid"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type ListFoodRes struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	sharedModel.DateDto
}
