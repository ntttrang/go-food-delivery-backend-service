package usergormmysql

import (
	"context"
	"time"

	usermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/user/model"
	shareinfras "github.com/ntttrang/go-food-delivery-backend-service/shared/infras"
	"github.com/pkg/errors"
)

type UserDeviceTokenRepo struct {
	dbCtx shareinfras.IDbContext
}

func NewUserDeviceTokenRepo(dbCtx shareinfras.IDbContext) *UserDeviceTokenRepo {
	return &UserDeviceTokenRepo{dbCtx: dbCtx}
}

func (r *UserDeviceTokenRepo) Insert(ctx context.Context, deviceToken *usermodel.UserDeviceToken) error {
	db := r.dbCtx.GetMainConnection()
	if err := db.Table(deviceToken.TableName()).Create(deviceToken).Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (r *UserDeviceTokenRepo) FindByToken(ctx context.Context, token string) (*usermodel.UserDeviceToken, error) {
	db := r.dbCtx.GetMainConnection()
	var deviceToken usermodel.UserDeviceToken

	// Find token regardless of revocation status - let business logic handle validation
	if err := db.Table(deviceToken.TableName()).Where("token = ?", token).First(&deviceToken).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	return &deviceToken, nil
}

func (r *UserDeviceTokenRepo) RevokeByUserId(ctx context.Context, userId string) error {
	db := r.dbCtx.GetMainConnection()
	now := time.Now().UTC()

	if err := db.Table(usermodel.UserDeviceToken{}.TableName()).
		Where("user_id = ? AND is_revoked = false", userId).
		Updates(map[string]interface{}{
			"is_revoked": true,
			"updated_at": now,
		}).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *UserDeviceTokenRepo) RevokeByToken(ctx context.Context, token string) error {
	db := r.dbCtx.GetMainConnection()
	now := time.Now().UTC()

	if err := db.Table(usermodel.UserDeviceToken{}.TableName()).
		Where("token = ?", token).
		Updates(map[string]interface{}{
			"is_revoked": true,
			"updated_at": now,
		}).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}
