package foodgormmysql

import (
	"context"

	"github.com/google/uuid"
	foodmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/food/model"
	"github.com/pkg/errors"
)

func (r *FoodRatingRepo) Delete(ctx context.Context, commentId uuid.UUID) error {
	db := r.dbCtx.GetMainConnection()

	if err := db.Where("id = ? ", commentId).Delete(foodmodel.FoodRatings{}).Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}
