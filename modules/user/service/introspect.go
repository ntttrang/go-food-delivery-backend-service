package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	usermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/user/model"
	shareComponent "github.com/ntttrang/go-food-delivery-backend-service/shared/component"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
	"golang.org/x/net/context"
)

// Define DTOs & validate
type IntrospectReq struct {
	Token string `json:"token"`
}

func (c *IntrospectReq) Validate() error {
	if c.Token == "" {
		return errors.New("token is required")
	}

	return nil
}

type IntrospectRes struct {
	Id        uuid.UUID            `json:"id"`
	LastName  string               `json:"last_name"`
	FirstName string               `json:"first_name"`
	Role      usermodel.UserRole   `json:"role"`
	Type      usermodel.UserType   `json:"type"`
	Status    usermodel.UserStatus `json:"status"`
	sharedModel.DateDto
}

func (ir *IntrospectRes) GetRole() uuid.UUID {
	return uuid.MustParse(string(ir.Role))
}

func (ir *IntrospectRes) Subject() uuid.UUID {
	return ir.Id
}

// Initilize service
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

// Implement
func (hdl *IntrospectCommandHandler) Execute(ctx context.Context, req IntrospectReq) (*IntrospectRes, error) {
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

	var res IntrospectRes
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
	introspectRes, err := w.hdl.Execute(context.Background(), IntrospectReq{Token: token})

	if err != nil {
		return nil, err
	}

	return introspectRes, nil
}
