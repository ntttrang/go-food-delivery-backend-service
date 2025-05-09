package gormmysql

import (
	"context"

	"github.com/google/uuid"
	"github.com/ntttrang/go-food-delivery-backend-service/modules/payment/model"
	"github.com/pkg/errors"
)

// FindByUserID finds cards by user ID
func (r *CardRepo) FindByUserID(ctx context.Context, userID uuid.UUID) ([]model.Card, error) {
	db := r.dbCtx.GetMainConnection()

	var cards []model.Card
	if err := db.WithContext(ctx).Where("user_id = ?", userID).Find(&cards).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	return cards, nil
}
