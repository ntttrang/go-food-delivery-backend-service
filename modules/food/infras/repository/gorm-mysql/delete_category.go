package foodgormmysql

import (
	"context"

	"github.com/google/uuid"

	foodmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/food/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	"github.com/pkg/errors"
)

func (r *FoodRepo) Delete(ctx context.Context, id uuid.UUID) error {
	db := r.dbCtx.GetMainConnection()

	if err := db.Table(foodmodel.Food{}.TableName()).
		Where("id = ?", id).
		Update("status", string(datatype.StatusDeleted)).
		Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}
