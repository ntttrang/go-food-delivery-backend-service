package gormmysql

import (
	"context"

	"github.com/pkg/errors"

	mediamodel "github.com/ntttrang/go-food-delivery-backend-service/modules/media/model"
)

func (repo *ImageRepository) Insert(ctx context.Context, data *mediamodel.Image) error {
	db := repo.dbCtx.GetMainConnection()

	if err := db.Create(data).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}
