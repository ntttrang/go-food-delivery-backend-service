package restaurantgormmysql

import (
	"context"

	"github.com/google/uuid"
	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (r *RestaurantFoodRepo) FindByRestaurantId(ctx context.Context, id uuid.UUID) ([]restaurantmodel.RestaurantFood, error) {
	var restaurantFoods []restaurantmodel.RestaurantFood

	db := r.dbCtx.GetMainConnection()
	if err := db.Where("restaurant_id = ?", id.String()).Find(&restaurantFoods).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, restaurantmodel.ErrRestaurantNotFound
		}
		return nil, errors.WithStack(err)
	}

	return restaurantFoods, nil
}
