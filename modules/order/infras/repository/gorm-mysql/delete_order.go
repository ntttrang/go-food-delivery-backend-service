package ordergormmysql

import (
	"context"

	ordermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/order/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	"github.com/pkg/errors"
)

func (r *OrderRepo) Delete(ctx context.Context, id string) error {
	db := r.dbCtx.GetMainConnection()

	// Start a transaction
	tx := db.Begin()
	if err := tx.Error; err != nil {
		return errors.WithStack(err)
	}

	// Soft delete order (update status to DELETED)
	if err := tx.Model(&ordermodel.Order{}).Where("id = ?", id).Update("status", datatype.StatusDeleted).Error; err != nil {
		tx.Rollback()
		return errors.WithStack(err)
	}

	// Soft delete order tracking
	if err := tx.Model(&ordermodel.OrderTracking{}).Where("order_id = ?", id).Update("status", datatype.StatusDeleted).Error; err != nil {
		tx.Rollback()
		return errors.WithStack(err)
	}

	// Soft delete order details
	if err := tx.Model(&ordermodel.OrderDetail{}).Where("order_id = ?", id).Update("status", datatype.StatusDeleted).Error; err != nil {
		tx.Rollback()
		return errors.WithStack(err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}
