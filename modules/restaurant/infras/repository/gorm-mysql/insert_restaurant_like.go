package restaurantgormmysql

import (
	"context"

	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	"github.com/pkg/errors"
)

func (r *RestaurantLikeRepo) Insert(ctx context.Context, data restaurantmodel.RestaurantLike) error {
	db := r.dbCtx.GetMainConnection()
	if err := db.Create(&data).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}
