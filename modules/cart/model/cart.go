package cartmodel

import (
	"github.com/google/uuid"
	sharedmodel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

const (
	CartStatusActive    = "ACTIVE"
	CartStatusUpdated   = "UPDATED"   // When Frontend update quantity
	CartStatusProcessed = "PROCESSED" // Auto updated by Backend. All items go to Order
)

type Cart struct {
	ID           uuid.UUID `gorm:"column:id;" json:"id"` // Can be duplicate
	UserID       uuid.UUID `gorm:"column:user_id" json:"userId"`
	FoodID       uuid.UUID `gorm:"column:food_id" json:"foodId"`
	RestaurantId uuid.UUID `gorm:"column:restaurant_id" json:"restaurantId"`
	Quantity     int       `gorm:"column:quantity" json:"quantity"`
	Status       string    `gorm:"column:status" json:"status"`
	DropOffLat   float64   `gorm:"column:dropoff_lat" json:"-"`
	DropOffLng   float64   `gorm:"column:dropoff_lng" json:"-"`
	sharedmodel.DateDto

	//ItemQuantity int64 `gorm:"item_quantity" json:"-"`
}

func (Cart) TableName() string {
	return "carts"
}

func (c *Cart) Validate() error {
	if c.UserID == uuid.Nil {
		return ErrUserIdRequired
	}

	if c.FoodID == uuid.Nil {
		return ErrFoodIdRequired
	}

	if c.Quantity <= 0 {
		return ErrQuantityInvalid
	}

	return nil
}
