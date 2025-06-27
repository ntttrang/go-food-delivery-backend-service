package cartgormmysql

import (
	"context"

	"github.com/google/uuid"
	cartmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/cart/model"
	"github.com/pkg/errors"
)

func (r *CartRepo) Delete(ctx context.Context, userId, foodId uuid.UUID) error {
	db := r.dbCtx.GetMainConnection()

	if err := db.Where("user_id = ? AND food_id = ?", userId, foodId).
		Delete(cartmodel.Cart{}).
		Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}
