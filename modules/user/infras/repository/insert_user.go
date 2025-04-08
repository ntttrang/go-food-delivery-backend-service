package userrepository

import (
	"context"

	usermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/user/model"
	"github.com/pkg/errors"
)

func (r *UserRepo) Insert(ctx context.Context, user *usermodel.User) error {
	if err := r.db.Table(user.TableName()).Create(user).Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}
