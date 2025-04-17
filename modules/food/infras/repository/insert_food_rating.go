package repository

import (
	"context"

	foodmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/food/model"
	"github.com/pkg/errors"
)

func (r *FoodRatingRepo) Insert(ctx context.Context, req *foodmodel.FoodCommentCreateReq) error {
	db := r.dbCtx.GetMainConnection()
	if err := db.Table(req.TableName()).Create(&req).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}
