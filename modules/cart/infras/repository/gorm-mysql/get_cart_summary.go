package cartgormmysql

import (
	"context"

	"github.com/google/uuid"
	cartmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/cart/model"
	"github.com/pkg/errors"
)

type CartSummaryData struct {
	CartID       uuid.UUID `json:"cartId"`
	UserID       uuid.UUID `json:"userId"`
	RestaurantID uuid.UUID `json:"restaurantId"`
	FoodId       uuid.UUID `json:"foodId"`
	Price        float64   `json:"price"`
	Quantity     int       `json:"quantity"`
	ItemCount    int       `json:"itemCount"`
	TotalPrice   float64   `json:"totalPrice"`
	Status       string    `json:"status"`
	CreatedAt    string    `json:"createdAt"`
	UpdatedAt    string    `json:"updatedAt"`
}

// GetCartSummaryByCartID gets cart summary by cart ID and user ID
func (r *CartRepo) GetCartSummaryByCartID(ctx context.Context, cartID uuid.UUID, userID uuid.UUID) (*CartSummaryData, error) {
	var result struct {
		CartID       uuid.UUID `gorm:"column:id"`
		UserID       uuid.UUID `gorm:"column:user_id"`
		RestaurantID uuid.UUID `gorm:"column:restaurant_id"`
		FoodId       uuid.UUID `json:"foodId"`
		Quantity     int       `gorm:"column:quantity"`
		Price        float64   `gorm:"column:price"`
		ItemCount    int       `gorm:"column:item_count"`
		TotalPrice   float64   `gorm:"column:total_price"`
		Status       string    `gorm:"column:status"`
		CreatedAt    string    `gorm:"column:created_at"`
		UpdatedAt    string    `gorm:"column:updated_at"`
	}

	db := r.dbCtx.GetMainConnection().Table(cartmodel.Cart{}.TableName())

	// Get cart summary with aggregated data
	// Note: This is a simplified implementation. In a real system, you might have:
	// 1. A separate carts table with cart_items
	// 2. Price information from food items
	// 3. More complex aggregation logic
	err := db.Select(`
		id,
		user_id,
		restaurant_id,
		food_id,
		quantity,
		0.0 as price,
		0 as item_count,
		0.0 as total_price,
		status,
		created_at,
		updated_at
	`).
		Where(" id = ? AND user_id = ? AND status != ?", cartID, userID, cartmodel.CartStatusProcessed).
		First(&result).Error

	if err != nil {
		if errors.Is(err, cartmodel.ErrCartNotFound) {
			return nil, cartmodel.ErrCartNotFound
		}
		return nil, errors.WithStack(err)
	}

	// Convert to service data structure
	return &CartSummaryData{
		CartID:       result.CartID,
		UserID:       result.UserID,
		RestaurantID: result.RestaurantID,
		FoodId:       result.FoodId,
		Quantity:     result.Quantity,
		Price:        result.Price,
		ItemCount:    result.ItemCount,
		TotalPrice:   result.TotalPrice,
		Status:       result.Status,
		CreatedAt:    result.CreatedAt,
		UpdatedAt:    result.UpdatedAt,
	}, nil
}
