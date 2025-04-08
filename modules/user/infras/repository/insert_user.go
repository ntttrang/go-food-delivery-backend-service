package userrepository

import (
	"context"

	usermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/user/model"
)

func (r *UserRepo) Insert(ctx context.Context, user *usermodel.User) error {
	if err := r.db.Table(user.TableName()).Create(user).Error; err != nil {
		return err
	}
	return nil
}
