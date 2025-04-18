package gormmysql

import (
	"context"

	mediamodel "github.com/ntttrang/go-food-delivery-backend-service/modules/media/model"

	"github.com/pkg/errors"
)

func (repo *ImageRepository) Insert(ctx context.Context, data *mediamodel.Image) error {
	db := repo.dbCtx.GetMainConnection()

	if err := db.Create(data).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}
