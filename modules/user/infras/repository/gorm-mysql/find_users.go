package usergormmysql

import (
	"context"

	usermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/user/model"
	service "github.com/ntttrang/go-food-delivery-backend-service/modules/user/service"
	"github.com/pkg/errors"
)

func (r *UserRepo) FindUsers(ctx context.Context, req service.UserListReq) ([]service.UserSearchResDto, int64, error) {
	db := r.dbCtx.GetMainConnection().Table(usermodel.User{}.TableName()).Select("id", "first_name", "last_name", "role", "email", "phone", "created_at", "updated_at") // Use field name ( Struct) or gorm name is OK
	if req.Name != "" {
		db = db.Where("first_name LIKE ? OR last_name LIKE ?", `%`+req.Name+`%`)
	}

	if req.Email != "" {
		db = db.Where("email = ?", req.Email)
	}

	if req.Phone != "" {
		db = db.Where("phone = ?", req.Phone)
	}

	if req.Role != "" {
		db = db.Where("role = ?", req.Role)
	}

	sortStr := "created_at DESC"
	if req.SortBy != "" {
		sortStr = req.SortBy + " " + req.Direction
	}

	var result []service.UserSearchResDto
	var total int64
	if err := db.Count(&total).Offset((req.Page - 1) * req.Limit).Limit(req.Limit).Order(sortStr).Find(&result).Error; err != nil {
		return nil, 0, errors.WithStack(err)
	}
	return result, total, nil
}
