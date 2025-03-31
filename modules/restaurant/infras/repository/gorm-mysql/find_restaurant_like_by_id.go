package restaurantgormmysql

import (
	"context"

	"github.com/google/uuid"
	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	"gorm.io/gorm"
)

func (repo *RestaurantLikeRepo) FindById(ctx context.Context, restaurantId uuid.UUID, userId uuid.UUID) (*restaurantmodel.RestaurantLike, error) {
	var restaurantLike restaurantmodel.RestaurantLike

	if err := repo.db.Where("restaurant_id = ? AND user_id = ?", restaurantId, userId).First(&restaurantLike).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, restaurantmodel.ErrRestaurantNotFound
		}
		return nil, err
	}

	return &restaurantLike, nil
}
