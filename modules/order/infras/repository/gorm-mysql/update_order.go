package ordergormmysql

import (
	"context"

	ordermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/order/model"
	"github.com/pkg/errors"
)

func (r *OrderRepo) Update(ctx context.Context, order *ordermodel.Order, tracking *ordermodel.OrderTracking) error {
	db := r.dbCtx.GetMainConnection()

	// Start a transaction
	tx := db.Begin()
	if err := tx.Error; err != nil {
		return errors.WithStack(err)
	}

	// Update order
	if err := tx.Save(order).Error; err != nil {
		tx.Rollback()
		return errors.WithStack(err)
	}

	// Update order tracking
	if err := tx.Save(tracking).Error; err != nil {
		tx.Rollback()
		return errors.WithStack(err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}
