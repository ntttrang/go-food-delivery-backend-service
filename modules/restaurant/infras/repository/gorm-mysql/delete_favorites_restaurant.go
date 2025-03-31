package restaurantgormmysql

import (
	"context"

	"github.com/google/uuid"
	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
)

func (repo *RestaurantLikeRepo) Delete(ctx context.Context, restaurantId uuid.UUID, userId uuid.UUID) error {
	return repo.db.Where("restaurant_id = ? AND user_id = ?", restaurantId, userId).Delete(restaurantmodel.RestaurantLike{}).Error
}
