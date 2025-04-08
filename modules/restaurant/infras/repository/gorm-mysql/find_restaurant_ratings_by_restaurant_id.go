package restaurantgormmysql

import (
	"context"

	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (repo *RestaurantRatingRepo) FindByRestaurantId(ctx context.Context, restaurantId string) ([]restaurantmodel.RestaurantRating, error) {
	var restaurantRating []restaurantmodel.RestaurantRating

	if err := repo.db.Where("restaurant_id = ?", restaurantId).Find(&restaurantRating).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, restaurantmodel.ErrRestaurantRatingNotFound
		}
		return nil, errors.WithStack(err)
	}

	return restaurantRating, nil
}
