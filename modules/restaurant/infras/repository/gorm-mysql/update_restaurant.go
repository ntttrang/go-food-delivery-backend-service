package restaurantgormmysql

import (
	"context"

	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	restaurantservice "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/service"
	"github.com/pkg/errors"
)

func (r *RestaurantRepo) Update(ctx context.Context, req restaurantservice.RestaurantUpdateReq) error {
	db := r.dbCtx.GetMainConnection()
	if err := db.Table(restaurantmodel.Restaurant{}.TableName()).Where("id = ?", req.Id).Updates(req.Dto).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}
