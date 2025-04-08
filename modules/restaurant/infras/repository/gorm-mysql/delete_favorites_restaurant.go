package restaurantgormmysql

import (
	"context"

	"github.com/google/uuid"
	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	"github.com/pkg/errors"
)

func (repo *RestaurantLikeRepo) Delete(ctx context.Context, restaurantId uuid.UUID, userId uuid.UUID) error {
	if err := repo.db.Where("restaurant_id = ? AND user_id = ?", restaurantId, userId).Delete(restaurantmodel.RestaurantLike{}).Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}
