package cartgormmysql

import (
	"context"

	"github.com/google/uuid"
	cartmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/cart/model"
	"github.com/pkg/errors"
)

func (r *CartRepo) UpdateStatusById(ctx context.Context, id uuid.UUID, status string) error {
	db := r.dbCtx.GetMainConnection()

	if err := db.Table(cartmodel.Cart{}.TableName()).Where("id = ?", id).Update("status", status).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}
