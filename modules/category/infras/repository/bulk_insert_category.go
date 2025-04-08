package repository

import (
	"context"

	categorymodel "github.com/ntttrang/go-food-delivery-backend-service/modules/category/model"
	"github.com/pkg/errors"
)

func (r *CategoryRepo) BulkInsert(ctx context.Context, data []categorymodel.Category) error {
	if err := r.db.Create(data).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}
