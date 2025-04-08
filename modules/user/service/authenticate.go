package service

import (
	"context"
	"errors"

	usermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/user/model"
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
		return nil, err
	}

	user, err := hdl.authRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if user != nil {
		if user.Status == usermodel.StatusDeleted || user.Status == usermodel.StatusBanned {
			return nil, errors.New("user is deleted or banned")
		}
	}

	// JWT
	token, err := hdl.tokenIssuer.IssueToken(ctx, user.Id.String())
	if err != nil {
		return nil, err
	}

	return &usermodel.AuthenticateRes{Token: token, ExpIn: hdl.tokenIssuer.ExpIn()}, nil
}
