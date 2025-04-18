package repository

import (
	"context"

	"github.com/google/uuid"
	foodmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/food/model"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
	"github.com/pkg/errors"
)

func (r *FoodRepo) Delete(ctx context.Context, id uuid.UUID) error {
	db := r.dbCtx.GetMainConnection()

	if err := db.Table(foodmodel.Food{}.TableName()).
		Where("id = ?", id).
		Update("status", sharedModel.StatusDelete).
		Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}
