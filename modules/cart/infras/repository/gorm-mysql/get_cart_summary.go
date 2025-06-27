package cartgormmysql

import (
	"context"

	"github.com/google/uuid"
	cartmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/cart/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	"github.com/pkg/errors"
)

type CartSummaryData struct {
	Id           uuid.UUID `json:"id"`
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
func (r *CartRepo) GetCartSummaryByCartID(ctx context.Context, cartID uuid.UUID, userID uuid.UUID) ([]CartSummaryData, error) {
	var summaries []CartSummaryData

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
		updated_at`).
		Where(" id = ? AND user_id = ? AND status != ?", cartID, userID, datatype.CartStatusProcessed).
		Find(&summaries).Error

	if err != nil {
		if errors.Is(err, cartmodel.ErrCartNotFound) {
			return nil, cartmodel.ErrCartNotFound
		}
		return nil, errors.WithStack(err)
	}

	return summaries, nil
}
