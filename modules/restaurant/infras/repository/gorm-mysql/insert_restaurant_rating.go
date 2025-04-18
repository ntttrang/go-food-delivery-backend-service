package restaurantgormmysql

import (
	"context"

	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	"github.com/pkg/errors"
)

func (r *RestaurantRatingRepo) Insert(ctx context.Context, req *restaurantmodel.RestaurantCommentCreateReq) error {
	db := r.dbCtx.GetMainConnection()
	if err := db.Table(req.TableName()).Create(&req).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}
