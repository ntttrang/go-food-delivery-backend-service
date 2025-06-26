package categorygormmysql

import (
	"context"

	categorymodel "github.com/ntttrang/go-food-delivery-backend-service/modules/category/model"
	"github.com/pkg/errors"
)

func (r *CategoryRepo) Insert(ctx context.Context, data *categorymodel.Category) error {
	db := r.dbCtx.GetMainConnection()
	if err := db.Create(data).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}
