package gormmysql

import (
	"context"

	"github.com/ntttrang/go-food-delivery-backend-service/modules/payment/model"
	"github.com/pkg/errors"
)

func (r *CardRepo) Create(ctx context.Context, card *model.Card) error {
	db := r.dbCtx.GetMainConnection()

	// Start a transaction
	tx := db.Begin()
	if err := tx.Error; err != nil {
		return errors.WithStack(err)
	}

	// Create the card
	if err := tx.WithContext(ctx).Create(card).Error; err != nil {
		tx.Rollback()
		return errors.WithStack(err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}
