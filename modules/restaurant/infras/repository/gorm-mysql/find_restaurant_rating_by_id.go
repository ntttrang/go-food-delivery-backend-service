package restaurantgormmysql

import (
	"context"

	"github.com/google/uuid"
	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (r *RestaurantRatingRepo) FindById(ctx context.Context, commentId uuid.UUID) (*restaurantmodel.RestaurantRating, error) {
	var restaurantRating restaurantmodel.RestaurantRating
	db := r.dbCtx.GetMainConnection()
	if err := db.Where("id = ?", commentId).First(&restaurantRating).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, restaurantmodel.ErrRestaurantNotFound
		}
		return nil, errors.WithStack(err)
	}

	return &restaurantRating, nil
}
