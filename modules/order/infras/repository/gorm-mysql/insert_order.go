package ordergormmysql

import (
	"context"

	ordermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/order/model"
	"github.com/pkg/errors"
)

func (r *OrderRepo) Insert(ctx context.Context, order *ordermodel.Order, orderTracking *ordermodel.OrderTracking, orderDetails []ordermodel.OrderDetail) error {
	db := r.dbCtx.GetMainConnection()

	// Start a transaction
	tx := db.Begin()
	if err := tx.Error; err != nil {
		return errors.WithStack(err)
	}

	// Insert order
	if err := tx.Create(order).Error; err != nil {
		tx.Rollback()
		return errors.WithStack(err)
	}

	// Insert order tracking
	if err := tx.Create(orderTracking).Error; err != nil {
		tx.Rollback()
		return errors.WithStack(err)
	}

	// Insert order details
	for _, detail := range orderDetails {
		if err := tx.Create(&detail).Error; err != nil {
			tx.Rollback()
			return errors.WithStack(err)
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}
