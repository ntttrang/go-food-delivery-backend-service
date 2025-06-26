package service

import (
	"context"
	"errors"
	"strings"

	usermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/user/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharemodel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

// Define DTOs & validate
type AuthenticateReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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
	Token string `json:"token"`
	ExpIn int    `json:"expIn"`
}

// Initilize service
type IAuthenticateRepo interface {
	FindByEmail(ctx context.Context, email string) (*usermodel.User, error)
}

type ITokenIssuer interface {
	IssueToken(ctx context.Context, userId string) (string, error)
	ExpIn() int
}

type AuthenticateCommandHandler struct {
	authRepo    IAuthenticateRepo
	tokenIssuer ITokenIssuer
}

func NewAuthenticateCommandHandler(authRepo IAuthenticateRepo, tokenIssuer ITokenIssuer) *AuthenticateCommandHandler {
	return &AuthenticateCommandHandler{
		authRepo:    authRepo,
		tokenIssuer: tokenIssuer,
	}
}

// Implement
func (hdl *AuthenticateCommandHandler) Execute(ctx context.Context, req AuthenticateReq) (*AuthenticateRes, error) {
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

	// JWT
	token, err := hdl.tokenIssuer.IssueToken(ctx, user.Id.String())
	if err != nil {
		return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	return &AuthenticateRes{Token: token, ExpIn: hdl.tokenIssuer.ExpIn()}, nil
}
