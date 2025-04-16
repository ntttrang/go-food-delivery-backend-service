package repository

import (
	"context"

	foodmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/food/model"
	"github.com/pkg/errors"
)

func (r *FoodRepo) Insert(ctx context.Context, data *foodmodel.Food) error {
	db := r.dbCtx.GetMainConnection()
	if err := db.Create(data).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}
