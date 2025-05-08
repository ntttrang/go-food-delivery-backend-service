package foodgormmysql

import (
	"context"

	"github.com/ntttrang/go-food-delivery-backend-service/modules/food/service"
	"github.com/pkg/errors"
)

func (r *FoodRatingRepo) Insert(ctx context.Context, req *service.FoodCommentCreateReq) error {
	db := r.dbCtx.GetMainConnection()
	if err := db.Table(req.TableName()).Create(&req).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}
