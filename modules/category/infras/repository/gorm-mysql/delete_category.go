package categorygormmysql

import (
	"context"

	"github.com/google/uuid"
	categorymodel "github.com/ntttrang/go-food-delivery-backend-service/modules/category/model"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
	"github.com/pkg/errors"
)

func (r *CategoryRepo) Delete(ctx context.Context, id uuid.UUID) error {
	db := r.dbCtx.GetMainConnection()

	if err := db.Table(categorymodel.Category{}.TableName()).
		Where("id = ?", id).
		Update("status", sharedModel.StatusDelete).
		Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}
