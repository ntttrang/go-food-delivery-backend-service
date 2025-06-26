package usergormmysql

import usermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/user/model"

func (r *UserAddressRepo) DeleteUserAddress(id string) error {
	db := r.dbCtx.GetMainConnection()
	if err := db.Where("id =?", id).Delete(&usermodel.UserAddress{}).Error; err != nil {
		return err
	}
	return nil
}
