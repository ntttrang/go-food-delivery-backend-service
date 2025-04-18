package restaurantmodel

import (
	"time"

	"github.com/google/uuid"
)

type MenuItemListRes struct {
	Items []MenuItemListDto `json:"items"`
}

type MenuItemListDto struct {
	FoodId       uuid.UUID  `json:"foodId"` // fooods.id
	Name         string     `json:"name"`   // foods.name
	Description  string     `json:"description"`
	ImageURL     string     `json:"imageUrl"`     // foods.images
	Price        float64    `json:"price"`        // foods.price
	Point        float64    `json:"point"`        // food_ratings.point
	CommentQty   int        `json:"commentQty"`   // food_ratings.comment
	CategoryId   uuid.UUID  `json:"categoryId"`   // food.category_id
	CategoryName string     `json:"categoryName"` // category.name
	RestaurantId uuid.UUID  `json:"restaurantId"`
	CreatedAt    *time.Time `json:"createdAt"`
	UpdatedAt    *time.Time `json:"updatedAt"`
}
