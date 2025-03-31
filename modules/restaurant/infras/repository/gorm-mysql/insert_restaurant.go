package restaurantgormmysql

import (
	"context"
	"fmt"

	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
)

func (repo *RestaurantRepo) Insert(ctx context.Context, restaurant restaurantmodel.Restaurant, restaurantFoods []restaurantmodel.RestaurantFood) error {
	tx := repo.db.Begin()
	if err := tx.Create(&restaurant).Error; err != nil {
		tx.Rollback()
		fmt.Println("create restaurant failed")
		return err
	}

	if len(restaurantFoods) > 0 {
		if err := tx.Create(&restaurantFoods).Error; err != nil {
			tx.Rollback()
			fmt.Println("create restaurant foods failed")
			return err
		}
	}

	return tx.Commit().Error
}
