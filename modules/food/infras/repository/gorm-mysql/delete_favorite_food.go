package foodgormmysql

import (
	"context"

	"github.com/google/uuid"
	foodmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/food/model"
	"github.com/pkg/errors"
)

func (r *FoodLikeRepo) Delete(ctx context.Context, foodId uuid.UUID, userId uuid.UUID) error {
	db := r.dbCtx.GetMainConnection()

	if err := db.Where("food_id = ? AND user_id = ?", foodId, userId).Delete(foodmodel.FoodLike{}).Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}
