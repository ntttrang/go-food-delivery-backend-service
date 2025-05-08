package usergormmysql

import (
	"context"

	usermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/user/model"
)

func (r *UserAddressRepo) InsertUserAddress(ctx context.Context, ua usermodel.UserAddress) error {
	db := r.dbCtx.GetMainConnection()
	result := db.Table(ua.TableName()).Create(&ua)
	return result.Error
}
