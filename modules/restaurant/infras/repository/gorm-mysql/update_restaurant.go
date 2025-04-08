package restaurantgormmysql

import (
	"context"

	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	"github.com/pkg/errors"
)

func (repo *RestaurantRepo) Update(ctx context.Context, req restaurantmodel.RestaurantUpdateReq) error {
	if err := repo.db.Table(restaurantmodel.Restaurant{}.TableName()).Where("id = ?", req.Id).Updates(req.Dto).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}
