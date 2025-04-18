package repository

import (
	"context"

	"github.com/google/uuid"
	categorymodel "github.com/ntttrang/go-food-delivery-backend-service/modules/category/model"
	"github.com/pkg/errors"
)

func (r *CategoryRepo) Update(ctx context.Context, id uuid.UUID, dto categorymodel.CategoryUpdateReq) error {
	db := r.dbCtx.GetMainConnection().Begin()

	if err := db.Table(dto.TableName()).Where("id = ?", id).Updates(dto).Error; err != nil {
		db.Rollback()
		return errors.WithStack(err)
	}

	if err := db.Commit().Error; err != nil {
		db.Rollback()
		return errors.WithStack(err)
	}

	return nil
}
