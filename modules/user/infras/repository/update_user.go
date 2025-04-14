package userrepository

import (
	"context"

	usermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/user/model"
	"github.com/pkg/errors"
)

func (r *UserRepo) Update(ctx context.Context, req *usermodel.UpdateUserReq) error {
	db := r.dbCtx.GetMainConnection()
	if err := db.Table(req.TableName()).Where("id = ?", req.Id).Updates(req).Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}
