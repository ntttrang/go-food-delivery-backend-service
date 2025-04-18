package service

import (
	"context"
	"errors"

	usermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/user/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

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

func (hdl *AuthenticateCommandHandler) Execute(ctx context.Context, req usermodel.AuthenticateReq) (*usermodel.AuthenticateRes, error) {
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
		if user.Status == usermodel.StatusDeleted || user.Status == usermodel.StatusBanned {
			return nil, datatype.ErrDeleted.WithError(usermodel.ErrUserDeletedOrBanned.Error())
		}
	}

	// JWT
	token, err := hdl.tokenIssuer.IssueToken(ctx, user.Id.String())
	if err != nil {
		return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	return &usermodel.AuthenticateRes{Token: token, ExpIn: hdl.tokenIssuer.ExpIn()}, nil
}
