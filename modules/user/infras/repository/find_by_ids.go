package userrepository

import (
	"context"

	"github.com/google/uuid"
	usermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/user/model"
)

func (r *UserRepo) FindByIds(ctx context.Context, ids []uuid.UUID) ([]usermodel.User, error) {
	var users []usermodel.User

	if err := r.dbCtx.GetMainConnection().Where("id IN (?)", ids).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}
