package restaurantgormmysql

import (
	"context"
	"fmt"

	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	"github.com/pkg/errors"
)

func (r *RestaurantRepo) Insert(ctx context.Context, restaurant restaurantmodel.Restaurant, restaurantFoods []restaurantmodel.RestaurantFood) error {
	tx := r.dbCtx.GetMainConnection().Begin()
	if err := tx.Create(&restaurant).Error; err != nil {
		tx.Rollback()
		fmt.Println("create restaurant failed")
		return errors.WithStack(err)
	}

	if len(restaurantFoods) > 0 {
		if err := tx.Create(&restaurantFoods).Error; err != nil {
			tx.Rollback()
			fmt.Println("create restaurant foods failed")
			return errors.WithStack(err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}
