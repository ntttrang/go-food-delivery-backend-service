package userrepository

import usermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/user/model"

func (r *UserAddressRepo) GetUserAddress(id string) (usermodel.UserAddress, error) {
	db := r.dbCtx.GetMainConnection()
	var ua usermodel.UserAddress
	result := db.Where("id =?", id).First(&ua)
	return ua, result.Error
}
