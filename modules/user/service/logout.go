package service

import (
	"context"

	usermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/user/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

type ILogoutRepo interface {
	FindByToken(ctx context.Context, token string) (*usermodel.UserDeviceToken, error)
	RevokeByToken(ctx context.Context, token string) error
}

type LogoutCommandHandler struct {
	deviceTokenRepo ILogoutRepo
}

func NewLogoutCommandHandler(deviceTokenRepo ILogoutRepo) *LogoutCommandHandler {
	return &LogoutCommandHandler{
		deviceTokenRepo: deviceTokenRepo,
	}
}

func (hdl *LogoutCommandHandler) Execute(ctx context.Context, refreshToken string) error {
	// Find the device token
	token, err := hdl.deviceTokenRepo.FindByToken(ctx, refreshToken)
	if err != nil {
		return datatype.ErrNotFound.WithError("Invalid refresh token")
	}

	// Check if token is valid
	if !token.IsValid() {
		return datatype.ErrUnauthorized.WithError("Refresh token is expired or revoked")
	}

	// Revoke the device token
	if err := hdl.deviceTokenRepo.RevokeByToken(ctx, refreshToken); err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	return nil
}
