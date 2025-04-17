package foodmodel

import (
	"github.com/google/uuid"
)

type FoodInfoDto struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Images      string    `json:"images"`
	Price       float64   `json:"price"`
	AvgPoint    float64   `json:"avgPoint"`
	CommentQty  int       `json:"commentQty"`
	CategoryId  uuid.UUID `json:"categoryId"`
}
