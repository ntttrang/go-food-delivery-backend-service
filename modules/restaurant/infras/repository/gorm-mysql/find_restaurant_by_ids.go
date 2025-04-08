package restaurantgormmysql

import (
	"context"

	"github.com/google/uuid"
	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	"github.com/pkg/errors"
)

func (repo *RestaurantRepo) FindRestaurantByIds(ctx context.Context, ids []uuid.UUID) ([]restaurantmodel.Restaurant, error) {
	var restaurants []restaurantmodel.Restaurant

	if err := repo.db.Where("id IN (?)", ids).Find(&restaurants).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	return restaurants, nil
}
