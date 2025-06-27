package restaurantgormmysql

import (
	"context"

	"github.com/google/uuid"
	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	"github.com/pkg/errors"
)

func (r *RestaurantRatingRepo) Delete(ctx context.Context, commentId uuid.UUID) error {
	db := r.dbCtx.GetMainConnection()

	if err := db.Where("id = ? ", commentId).Delete(restaurantmodel.RestaurantRating{}).Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}
