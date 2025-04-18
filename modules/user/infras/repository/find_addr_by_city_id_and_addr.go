package userrepository

import (
	"context"

	usermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/user/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (r *UserAddressRepo) FindUsrAddrByCityIdAndAddr(ctx context.Context, cityId int, addr string) (*usermodel.UserAddress, error) {
	db := r.dbCtx.GetMainConnection()
	var ua usermodel.UserAddress
	if err := db.Table(usermodel.UserAddress{}.TableName()).Where("city_id = ? AND addr = ?", cityId, addr).First(&ua).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, usermodel.ErrUserAddrNotFound
		}
		return nil, errors.WithStack(err)
	}
	return &ua, nil
}
