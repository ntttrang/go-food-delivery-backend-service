package service

import (
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	usermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/user/model"
	shareComponent "github.com/ntttrang/go-food-delivery-backend-service/shared/component"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	"golang.org/x/net/context"
)

type IIntrospectRepo interface {
	FindById(ctx context.Context, id uuid.UUID) (*usermodel.User, error)
}

type IntrospectCommandHandler struct {
	jwtComp *shareComponent.JwtComp
	repo    IIntrospectRepo
}

func NewIntrospectCommandHandler(jwtComp *shareComponent.JwtComp, repo IIntrospectRepo) *IntrospectCommandHandler {
	return &IntrospectCommandHandler{
		jwtComp: jwtComp,
		repo:    repo,
	}
}

func (hdl *IntrospectCommandHandler) Execute(ctx context.Context, req usermodel.IntrospectReq) (*usermodel.IntrospectRes, error) {
	if err := req.Validate(); err != nil {
		return nil, datatype.ErrBadRequest.WithWrap(err).WithDebug(err.Error())
	}

	userId, err := hdl.jwtComp.Validate(req.Token)
	if err != nil {
		return nil, datatype.ErrUnauthorized.WithWrap(err)
	}

	user, err := hdl.repo.FindById(ctx, uuid.MustParse(userId))
	if err != nil {
		return nil, datatype.ErrUnauthorized.WithWrap(err)
	}

	if user.Status == usermodel.StatusBanned || user.Status == usermodel.StatusDeleted {
		return nil, datatype.ErrUnauthorized.WithDebug(usermodel.ErrUserDeletedOrBanned.Error())
	}

	var res usermodel.IntrospectRes
	err = copier.Copy(&res, &user)
	if err != nil {
		return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug("copier.Copy failed")
	}

	return &res, nil
}

// Implement of ITokenIntrospector
type IntrospectCmdHdlWrapper struct {
	hdl *IntrospectCommandHandler
}

func NewIntrospectCmdHdlWrapper(hdl *IntrospectCommandHandler) *IntrospectCmdHdlWrapper {
	return &IntrospectCmdHdlWrapper{hdl: hdl}
}

func (w *IntrospectCmdHdlWrapper) Validate(token string) (datatype.Requester, error) {
	introspectRes, err := w.hdl.Execute(context.Background(), usermodel.IntrospectReq{Token: token})

	if err != nil {
		return nil, err
	}

	return introspectRes, nil
}
