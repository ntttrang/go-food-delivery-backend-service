package restaurantgormmysql

import (
	"context"

	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
)

func (repo *RestaurantRepo) Update(ctx context.Context, req restaurantmodel.RestaurantUpdateReq) error {
	return repo.db.Table(restaurantmodel.Restaurant{}.TableName()).Where("id = ?", req.Id).Updates(req.Dto).Error
}
