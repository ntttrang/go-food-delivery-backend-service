package cartgormmysql

import (
	"context"

	"github.com/google/uuid"
	cartmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/cart/model"
	"github.com/pkg/errors"
)

// UpdateCartStatusByCartID updates cart status by cart ID
func (r *CartRepo) UpdateCartStatusByCartID(ctx context.Context, cartID uuid.UUID, status string) error {
	db := r.dbCtx.GetMainConnection().Table(cartmodel.Cart{}.TableName())

	// Update all cart items with the same cart ID (assuming cart ID represents a group of items)
	// In this simplified implementation, we'll update by user's cart items
	result := db.Where("id = ?", cartID).Update("status", status)

	if result.Error != nil {
		return errors.WithStack(result.Error)
	}

	// Check if any rows were affected
	if result.RowsAffected == 0 {
		return cartmodel.ErrCartNotFound
	}

	return nil
}
