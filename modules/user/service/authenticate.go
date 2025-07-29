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
	sharemodel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
	"golang.org/x/crypto/bcrypt"
)

// Define DTOs & validate
type AuthenticateReq struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	OS           string `json:"os,omitempty"`           // Optional: can be set by client
	IsProduction *bool  `json:"isProduction,omitempty"` // Optional: can be set by client
}

func (r *AuthenticateReq) Validate() error {
	r.Email = strings.TrimSpace(r.Email)
	r.Password = strings.TrimSpace(r.Password)

	if r.Email == "" {
		return usermodel.ErrEmailRequired
	}

	if r.Password == "" {
		return usermodel.ErrPasswordInvalid
	}

	if !sharemodel.ValidateEmail(r.Email) {
		return usermodel.ErrEmailInvalid
	}

	if len(r.Password) <= 6 {
		return usermodel.ErrPasswordInvalid
	}

	return nil
}

type AuthenticateRes struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpIn        int    `json:"expIn"`
}

// Initilize service
type IAuthenticateRepo interface {
	FindByEmail(ctx context.Context, email string) (*usermodel.User, error)
}

type IUserDeviceTokenRepo interface {
	Insert(ctx context.Context, deviceToken *usermodel.UserDeviceToken) error
	FindByToken(ctx context.Context, token string) (*usermodel.UserDeviceToken, error)
	RevokeByUserId(ctx context.Context, userId string) error
	RevokeByToken(ctx context.Context, token string) error
}

type ITokenIssuer interface {
	IssueToken(ctx context.Context, userId string) (string, error)
	ExpIn() int
}

type IRefreshTokenGenerator interface {
	GenerateRefreshToken() (string, error)
	RefreshTokenExpiry() time.Duration
}

type AuthenticateCommandHandler struct {
	authRepo        IAuthenticateRepo
	deviceTokenRepo IUserDeviceTokenRepo
	tokenIssuer     ITokenIssuer
	tokenGenerator  IRefreshTokenGenerator
}

func NewAuthenticateCommandHandler(
	authRepo IAuthenticateRepo,
	deviceTokenRepo IUserDeviceTokenRepo,
	tokenIssuer ITokenIssuer,
	tokenGenerator IRefreshTokenGenerator,
) *AuthenticateCommandHandler {
	return &AuthenticateCommandHandler{
		authRepo:        authRepo,
		deviceTokenRepo: deviceTokenRepo,
		tokenIssuer:     tokenIssuer,
		tokenGenerator:  tokenGenerator,
	}
}

// Implement
func (hdl *AuthenticateCommandHandler) Execute(ctx context.Context, req AuthenticateReq, userAgent string) (*AuthenticateRes, error) {
	if err := req.Validate(); err != nil {
		return nil, datatype.ErrBadRequest.WithWrap(err).WithDebug(err.Error())
	}

	user, err := hdl.authRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, usermodel.ErrUserNotFound) {
			return nil, datatype.ErrNotFound.WithDebug(usermodel.ErrUserNotFound.Error())
		}
		return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	if user != nil {
		if user.Status == datatype.StatusDeleted || user.Status == datatype.StatusBanned {
			return nil, datatype.ErrDeleted.WithError(usermodel.ErrUserDeletedOrBanned.Error())
		}
	}

	// Verify password against stored hash
	// Password is hashed using format: salt.password
	saltPass := fmt.Sprintf("%s.%s", user.Salt, req.Password)
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(saltPass)); err != nil {
		// Return the same error as user not found to prevent user enumeration
		return nil, datatype.ErrNotFound.WithDebug("invalid credentials")
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
