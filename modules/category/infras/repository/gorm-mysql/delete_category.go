package categorygormmysql

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	categorymodel "github.com/ntttrang/go-food-delivery-backend-service/modules/category/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

func (r *CategoryRepo) Delete(ctx context.Context, id uuid.UUID) error {
	db := r.dbCtx.GetMainConnection()

	if err := db.Table(categorymodel.Category{}.TableName()).
		Where("id = ?", id).
		Update("status", datatype.StatusDeleted).
		Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}
