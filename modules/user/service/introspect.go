package service

import (
	"errors"

	"github.com/google/uuid"
	usermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/user/model"
	sharecomponent "github.com/ntttrang/go-food-delivery-backend-service/shared/component"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharemodel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
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
	Id        uuid.UUID           `json:"id"`
	LastName  string              `json:"last_name"`
	FirstName string              `json:"first_name"`
	Role      datatype.UserRole   `json:"role"`
	Type      datatype.UserType   `json:"type"`
	Status    datatype.UserStatus `json:"status"`
	sharemodel.DateDto
}

func (ir *IntrospectRes) GetRole() string {
	return string(ir.Role)
}

func (ir *IntrospectRes) Subject() uuid.UUID {
	return ir.Id
}

// Initilize service
type IIntrospectRepo interface {
	FindById(ctx context.Context, id uuid.UUID) (*usermodel.User, error)
}

type IntrospectCommandHandler struct {
	jwtComp *sharecomponent.JwtComp
	repo    IIntrospectRepo
}

func NewIntrospectCommandHandler(jwtComp *sharecomponent.JwtComp, repo IIntrospectRepo) *IntrospectCommandHandler {
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

	if user.Status == datatype.StatusBanned || user.Status == datatype.StatusDeleted {
		return nil, datatype.ErrUnauthorized.WithDebug(usermodel.ErrUserDeletedOrBanned.Error())
	}

	var res = IntrospectRes{
		Id:        user.Id,
		LastName:  user.LastName,
		FirstName: user.FirstName,
		Role:      user.Role,
		Type:      user.Type,
		Status:    user.Status,
		DateDto: sharemodel.DateDto{
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
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
