package usergormmysql

import usermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/user/model"

func (r *UserAddressRepo) UpdateUserAddress(ua usermodel.UserAddress) error {
	db := r.dbCtx.GetMainConnection()
	result := db.Model(&ua).Updates(ua)
	return result.Error
}
