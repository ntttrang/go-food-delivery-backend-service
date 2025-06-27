package restaurantgormmysql

import (
	"context"

	"github.com/google/uuid"
	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (r *RestaurantLikeRepo) FindByUserIdAndRestaurantId(ctx context.Context, userId, restaurantId uuid.UUID) (*restaurantmodel.RestaurantLike, error) {
	var restaurantLike restaurantmodel.RestaurantLike
	db := r.dbCtx.GetMainConnection()
	if err := db.Where("restaurant_id = ? AND user_id = ?", restaurantId, userId).First(&restaurantLike).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, restaurantmodel.ErrRestaurantLikeNotFound
		}
		return nil, errors.WithStack(err)
	}

	return &restaurantLike, nil
}
