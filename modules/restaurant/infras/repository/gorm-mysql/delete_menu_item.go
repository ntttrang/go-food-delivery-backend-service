package restaurantgormmysql

import (
	"context"

	rpcclient "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/infras/repository/rpc-client"
	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	restaurantservice "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/service"
	"github.com/pkg/errors"
)

func (r *RestaurantFoodRepo) Delete(ctx context.Context, req *restaurantservice.MenuItemCreateReq) error {

	tx := r.dbCtx.GetMainConnection().Begin()

	if err := tx.Where("restaurant_id = ? AND food_id = ?", req.RestaurantId, req.FoodId).Delete(restaurantmodel.RestaurantFood{}).Error; err != nil {
		tx.Rollback()
		return errors.WithStack(err)
	}

	updateReq := rpcclient.FoodUpdateReq{
		RestaurantId: "",
		CategoryId:   "",
		FoodId:       req.FoodId.String(),
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
