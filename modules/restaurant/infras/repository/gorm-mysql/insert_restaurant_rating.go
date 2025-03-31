package restaurantgormmysql

import (
	"context"

	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
)

func (repo *RestaurantRatingRepo) Insert(ctx context.Context, req *restaurantmodel.RestaurantCommentCreateReq) error {
	if err := repo.db.Table(req.TableName()).Create(&req).Error; err != nil {
		return err
	}

	return nil
}
