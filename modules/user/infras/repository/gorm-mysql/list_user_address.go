package usergormmysql

import (
	"context"

	"github.com/google/uuid"
	service "github.com/ntttrang/go-food-delivery-backend-service/modules/user/service"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

func (r *UserAddressRepo) ListUserAddresses(ctx context.Context, userId uuid.UUID) ([]service.UserAddrSearchResDto, error) {
	db := r.dbCtx.GetMainConnection()
	var userAddresses []service.UserAddrSearchResDto
	if err := db.Raw(`SELECT 
						ua.id,
						ua.user_id,
						ua.city_id,
						c.title AS city_name,
						ua.addr,
						ua.lat,
						ua.lng
					FROM user_addresses ua
					INNER JOIN cities c
					ON ua.city_id = c.id
					WHERE ua.user_id  = ? AND ua.status = ? AND c.status = 1
					ORDER BY ua.created_at DESC`, userId, datatype.StatusActive).Find(&userAddresses).Error; err != nil {
		return nil, err
	}
	return userAddresses, nil
}
