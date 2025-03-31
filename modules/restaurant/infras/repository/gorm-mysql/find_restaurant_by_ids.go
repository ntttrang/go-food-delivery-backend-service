package restaurantgormmysql

import (
	"context"

	"github.com/google/uuid"
	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
)

func (repo *RestaurantRepo) FindRestaurantByIds(ctx context.Context, ids []uuid.UUID) ([]restaurantmodel.Restaurant, error) {
	var restaurants []restaurantmodel.Restaurant

	if err := repo.db.Where("id IN (?)", ids).Find(&restaurants).Error; err != nil {
		return nil, err
	}

	return restaurants, nil
}
