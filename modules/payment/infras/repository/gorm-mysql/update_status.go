package gormmysql

import (
	"context"

	"github.com/google/uuid"
	"github.com/ntttrang/go-food-delivery-backend-service/modules/payment/model"
	"github.com/pkg/errors"
)

func (r *CardRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	db := r.dbCtx.GetMainConnection()

	// Start a transaction
	tx := db.Begin()
	if err := tx.Error; err != nil {
		return errors.WithStack(err)
	}

	// Update the card status
	if err := tx.WithContext(ctx).Model(&model.Card{}).Where("id = ?", id).Update("status", status).Error; err != nil {
		tx.Rollback()
		return errors.WithStack(err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}
