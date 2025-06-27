package gormmysql

import (
	"context"

	"github.com/google/uuid"
	"github.com/ntttrang/go-food-delivery-backend-service/modules/payment/model"
	"github.com/pkg/errors"
)

func (r *CardRepo) FindByID(ctx context.Context, id uuid.UUID) (*model.Card, error) {
	db := r.dbCtx.GetMainConnection()

	var card model.Card
	if err := db.WithContext(ctx).Where("id = ?", id).First(&card).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	return &card, nil
}
