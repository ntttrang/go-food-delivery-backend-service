package restaurantgormmysql

import (
	"context"

	"github.com/google/uuid"
	rpcclient "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/infras/repository/rpc-client"
	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	"github.com/pkg/errors"
)

func (r *RestaurantFoodRepo) Insert(ctx context.Context, restaurantFood *restaurantmodel.RestaurantFood, categoryId uuid.UUID) error {

	tx := r.dbCtx.GetMainConnection().Begin()

	if err := tx.Table(restaurantmodel.RestaurantFood{}.TableName()).Create(&restaurantFood).Error; err != nil {
		tx.Rollback()
		return errors.WithStack(err)
	}

	updateReq := rpcclient.FoodUpdateReq{
		RestaurantId: restaurantFood.RestaurantId.String(),
		CategoryId:   categoryId.String(),
		FoodId:       restaurantFood.FoodId.String(),
	}
	_, errR := r.foodRPCClient.UpdateFoods(ctx, updateReq)
	if errR != nil {
		tx.Rollback()
		return errors.WithStack(errR)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
