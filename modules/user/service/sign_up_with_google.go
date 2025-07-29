package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	usermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/user/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharedmodel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type ISignupGoogleRepo interface {
	FindByEmail(ctx context.Context, email string) (*usermodel.User, error)
	Insert(ctx context.Context, user *usermodel.User) error
}

type IGoogleTokenIssuer interface {
	IssueToken(ctx context.Context, userId string) (string, error)
	ExpIn() int
}

type IGoogleOAuth interface {
	GenerateState(ctx context.Context) string
	AuthCodeUrl(ctx context.Context, state string) string
	GetGGUserInfo(ctx context.Context, state string, code string) (*sharedmodel.GgUserInfo, error)
}

type SignUpGoogleCommandHandler struct {
	repo            ISignupGoogleRepo
	deviceTokenRepo IUserDeviceTokenRepo
	tokenIssuer     IGoogleTokenIssuer
	tokenGenerator  IRefreshTokenGenerator
	ggOAuth         IGoogleOAuth
}

func NewSignUpGoogleCommandHandler(
	repo ISignupGoogleRepo,
	deviceTokenRepo IUserDeviceTokenRepo,
	tokenIssuer IGoogleTokenIssuer,
	tokenGenerator IRefreshTokenGenerator,
	ggOAuth IGoogleOAuth,
) *SignUpGoogleCommandHandler {
	return &SignUpGoogleCommandHandler{
		repo:            repo,
		deviceTokenRepo: deviceTokenRepo,
		tokenIssuer:     tokenIssuer,
		tokenGenerator:  tokenGenerator,
		ggOAuth:         ggOAuth,
	}
}

func (hdl *SignUpGoogleCommandHandler) GetAuthCodeUrl(ctx context.Context) string {
	state := hdl.ggOAuth.GenerateState(ctx)
	url := hdl.ggOAuth.AuthCodeUrl(ctx, state)
	fmt.Printf("state = %s, url = %s \n", state, url)
	return url
}

func (hdl *SignUpGoogleCommandHandler) AuthenticateByGoogle(ctx context.Context, state string, code string) (*AuthenticateRes, error) {
	ggUserInfo, err := hdl.ggOAuth.GetGGUserInfo(ctx, state, code)
	if err != nil {
		return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	email := ggUserInfo.Email
	user, err := hdl.repo.FindByEmail(ctx, email)
	if err != nil {
		if !errors.Is(err, usermodel.ErrUserNotFound) {
			return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
		}
	}

	if user != nil {
		if user.Status == datatype.StatusDeleted || user.Status == datatype.StatusBanned {
			return nil, datatype.ErrDeleted.WithError(usermodel.ErrUserDeletedOrBanned.Error())
		}
	}

	if errors.Is(err, usermodel.ErrUserNotFound) {
		user = &usermodel.User{}
		// Create user
		user.Id, _ = uuid.NewV7()
		user.Email = email
		user.Status = datatype.StatusActive // Always set Active Status when insert
		nameParts := strings.Split(ggUserInfo.Name, " ")

		firstName := ""
		lastName := ""

		if len(nameParts) > 0 {
			firstName = nameParts[0]

			if len(nameParts) > 1 {
				lastName = strings.Join(nameParts[1:], " ")
			}
		}
		user.FirstName = firstName
		user.LastName = lastName
		user.GgId = ggUserInfo.GgId
		user.Role = datatype.RoleUser
		user.Type = datatype.TypeGmail
		now := time.Now().UTC()
		user.CreatedAt = &now
		user.UpdatedAt = &now
		if err := hdl.repo.Insert(ctx, user); err != nil {
			return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
		}
	}

	// Revoke existing device tokens for this user
	if err := hdl.deviceTokenRepo.RevokeByUserId(ctx, user.Id.String()); err != nil {
		return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	// Generate access token
	accessToken, err := hdl.tokenIssuer.IssueToken(ctx, user.Id.String())
	if err != nil {
		return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	// Generate refresh token
	refreshToken, err := hdl.tokenGenerator.GenerateRefreshToken()
	if err != nil {
		return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	// Save device token to database
	deviceTokenModel := &usermodel.UserDeviceToken{
		Id:           uuid.New(),
		UserId:       user.Id,
		Token:        refreshToken,
		ExpiresAt:    time.Now().UTC().Add(hdl.tokenGenerator.RefreshTokenExpiry()),
		IsRevoked:    false,
		IsProduction: true, // Default to production
		OS:           "",   // Will be set from request headers in future
	}
	now := time.Now().UTC()
	deviceTokenModel.CreatedAt = &now
	deviceTokenModel.UpdatedAt = &now

	if err := hdl.deviceTokenRepo.Insert(ctx, deviceTokenModel); err != nil {
		return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	return &AuthenticateRes{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpIn:        hdl.tokenIssuer.ExpIn(),
	}, nil
}
