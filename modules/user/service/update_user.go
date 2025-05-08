package service

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	usermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/user/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

// Define DTOs & validate
type UpdateUserReq struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Phone     string `json:"phone"`
	Avatar    string `json:"avatar"`
	Status    string `json:"status"`

	Id uuid.UUID `json:"_"`
}

func (UpdateUserReq) TableName() string {
	return usermodel.User{}.TableName()
}

func (r *UpdateUserReq) Validate() error {
	r.FirstName = strings.TrimSpace(r.FirstName)
	r.LastName = strings.TrimSpace(r.LastName)
	r.Phone = strings.TrimSpace(r.Phone)

	if r.Id.String() == "" {
		return usermodel.ErrIdRequired
	}

	if r.Phone != "" && !sharedModel.ValidatePhoneNumber(r.Phone) {
		return usermodel.ErrInvalidPhoneNumber
	}

	return nil
}

// Initilize service
type IupdateRepo interface {
	FindById(ctx context.Context, id uuid.UUID) (*usermodel.User, error)
	Update(ctx context.Context, req *UpdateUserReq) error
}

type UpdateCommandHandler struct {
	userRepo IupdateRepo
}

func NewUpdateCommandHandler(userRepo IupdateRepo) *UpdateCommandHandler {
	return &UpdateCommandHandler{userRepo: userRepo}
}

// Implement
func (hdl *UpdateCommandHandler) Execute(ctx context.Context, req UpdateUserReq) error {
	if err := req.Validate(); err != nil {
		return datatype.ErrBadRequest.WithWrap(err).WithDebug(err.Error())
	}

	existUser, err := hdl.userRepo.FindById(ctx, req.Id)

	if err != nil {
		if errors.Is(err, usermodel.ErrUserNotFound) {
			return datatype.ErrNotFound.WithDebug(usermodel.ErrUserNotFound.Error())
		}

		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	if existUser.Status == sharedModel.StatusDelete {
		return datatype.ErrNotFound.WithError(usermodel.ErrUserDeletedOrBanned.Error())
	}

	if err := hdl.userRepo.Update(ctx, &req); err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	return nil
}
