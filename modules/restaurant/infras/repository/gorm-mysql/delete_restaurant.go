package restaurantgormmysql

import (
	"context"

	"github.com/google/uuid"
	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	"github.com/pkg/errors"
)

func (r *RestaurantRepo) Delete(ctx context.Context, id uuid.UUID) error {
	db := r.dbCtx.GetMainConnection()
	if err := db.Table(restaurantmodel.Restaurant{}.TableName()).Where("id = ?", id).Update("status", string(datatype.StatusDeleted)).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}
