package usergormmysql

import (
	"context"

	"github.com/google/uuid"
	usermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/user/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (repo *UserRepo) FindByEmail(ctx context.Context, email string) (*usermodel.User, error) {
	return repo.FindByCondition(ctx, map[string]interface{}{"email": email})
}

func (repo *UserRepo) FindById(ctx context.Context, id uuid.UUID) (*usermodel.User, error) {
	return repo.FindByCondition(ctx, map[string]interface{}{"id": id})
}

func (r *UserRepo) FindByCondition(ctx context.Context, cond map[string]interface{}) (*usermodel.User, error) {
	var user usermodel.User
	db := r.dbCtx.GetMainConnection()
	if err := db.Table(user.TableName()).Where(cond).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, usermodel.ErrUserNotFound
		}

		return nil, errors.WithStack(err)
	}

	return &user, nil
}
