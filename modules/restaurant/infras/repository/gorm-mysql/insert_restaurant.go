package restaurantgormmysql

import (
	"context"
	"fmt"

	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	"github.com/pkg/errors"
)

func (r *RestaurantRepo) Insert(ctx context.Context, restaurant restaurantmodel.Restaurant) error {
	tx := r.dbCtx.GetMainConnection().Begin()
	if err := tx.Create(&restaurant).Error; err != nil {
		tx.Rollback()
		fmt.Println("create restaurant failed")
		return errors.WithStack(err)
	}

	if err := tx.Commit().Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}
