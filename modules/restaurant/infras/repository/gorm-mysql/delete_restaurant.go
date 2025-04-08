package restaurantgormmysql

import (
	"context"

	"github.com/google/uuid"
	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
	"github.com/pkg/errors"
)

func (repo *RestaurantRepo) Delete(ctx context.Context, id uuid.UUID) error {
	if err := repo.db.Table(restaurantmodel.Restaurant{}.TableName()).Where("id = ?", id).Update("status", sharedModel.StatusDelete).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}
