package restaurantgormmysql

import (
	"context"

	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/internal/model"
)

func (repo *RestaurantFoodRepo) BulkInsert(ctx context.Context, data []restaurantmodel.RestaurantFood) error {
	if err := repo.db.Create(data).Error; err != nil {
		return err
	}

	return nil
}
