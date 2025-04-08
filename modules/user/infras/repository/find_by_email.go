package userrepository

import (
	"context"

	usermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/user/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*usermodel.User, error) {
	var user *usermodel.User

	if err := r.db.Where("email = ?", email).Find(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, usermodel.ErrUserNotFound
		}

		return nil, errors.WithStack(err)
	}

	return user, nil
}
