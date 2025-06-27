package restaurantgormmysql

import (
	"context"

	"github.com/google/uuid"
	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	"github.com/pkg/errors"
)

func (r *RestaurantFoodRepo) FindFoodByRestaurantIds(ctx context.Context, ids []uuid.UUID) ([]restaurantmodel.RestaurantFood, error) {
	var restaurantfoods []restaurantmodel.RestaurantFood

	db := r.dbCtx.GetMainConnection()
	if err := db.Where("restaurant_id IN (?)", ids).Find(&restaurantfoods).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	return restaurantfoods, nil
}
