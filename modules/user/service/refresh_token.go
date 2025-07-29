package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	usermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/user/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

// Interfaces are defined in authenticate.go to avoid duplication

type RefreshTokenCommandHandler struct {
	deviceTokenRepo IUserDeviceTokenRepo
	tokenIssuer     ITokenIssuer
	tokenGenerator  IRefreshTokenGenerator
}

func NewRefreshTokenCommandHandler(
	deviceTokenRepo IUserDeviceTokenRepo,
	tokenIssuer ITokenIssuer,
	tokenGenerator IRefreshTokenGenerator,
) *RefreshTokenCommandHandler {
	return &RefreshTokenCommandHandler{
		deviceTokenRepo: deviceTokenRepo,
		tokenIssuer:     tokenIssuer,
		tokenGenerator:  tokenGenerator,
	}
}

func (hdl *RefreshTokenCommandHandler) Execute(ctx context.Context, refreshToken string) (*AuthenticateRes, error) {
	// Find the device token
	token, err := hdl.deviceTokenRepo.FindByToken(ctx, refreshToken)
	if err != nil {
		return nil, datatype.ErrNotFound.WithError("Invalid refresh token")
	}

	// Check if token is valid
	if !token.IsValid() {
		return nil, datatype.ErrUnauthorized.WithError("Refresh token is expired or revoked")
	}

	// Generate new access token
	accessToken, err := hdl.tokenIssuer.IssueToken(ctx, token.UserId.String())
	if err != nil {
		return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	// Generate new refresh token
	newRefreshToken, err := hdl.tokenGenerator.GenerateRefreshToken()
	if err != nil {
		return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	// Revoke old device token
	if err := hdl.deviceTokenRepo.RevokeByToken(ctx, refreshToken); err != nil {
		return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	// Save new device token
	newToken := &usermodel.UserDeviceToken{
		Id:           uuid.New(),
		UserId:       token.UserId,
		Token:        newRefreshToken,
		ExpiresAt:    time.Now().UTC().Add(hdl.tokenGenerator.RefreshTokenExpiry()),
		IsRevoked:    false,
		IsProduction: token.IsProduction, // Keep same environment
		OS:           token.OS,           // Keep same OS
	}
	now := time.Now().UTC()
	newToken.CreatedAt = &now
	newToken.UpdatedAt = &now

	if err := hdl.deviceTokenRepo.Insert(ctx, newToken); err != nil {
		return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	return &AuthenticateRes{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpIn:        hdl.tokenIssuer.ExpIn(),
	}, nil
}
