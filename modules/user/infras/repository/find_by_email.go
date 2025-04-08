package userrepository

import (
	"context"

	usermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/user/model"
)

func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*usermodel.User, error) {
	var user *usermodel.User

	if err := r.db.Where("email = ?", email).Find(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}
