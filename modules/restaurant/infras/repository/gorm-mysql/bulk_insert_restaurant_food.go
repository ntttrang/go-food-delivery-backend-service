package restaurantgormmysql

import (
	"context"

	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	"github.com/pkg/errors"
)

func (r *RestaurantFoodRepo) BulkInsert(ctx context.Context, data []restaurantmodel.RestaurantFood) error {
	db := r.dbCtx.GetMainConnection()

	if err := db.Create(data).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}
