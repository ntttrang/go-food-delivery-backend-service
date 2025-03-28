package restaurantgormmysql

import (
	"context"

	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
)

func (repo *RestaurantRepo) Insert(ctx context.Context, data restaurantmodel.Restaurant) error {
	if err := repo.db.Create(&data).Error; err != nil {
		return err
	}

	return nil
}
