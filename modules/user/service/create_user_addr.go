package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	usermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/user/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type ICreateUserAddrRepo interface {
	FindUsrAddrByCityIdAndAddr(ctx context.Context, cityId int, addr string) (*usermodel.UserAddress, error)
	InsertUserAddress(ctx context.Context, ua usermodel.UserAddress) error
}

type CreateUserAddrCommandHandler struct {
	userAddrRepo ICreateUserAddrRepo
}

func NewCreateUserAddrCommandHandler(userAddrRepo ICreateUserAddrRepo) *CreateUserAddrCommandHandler {
	return &CreateUserAddrCommandHandler{userAddrRepo: userAddrRepo}
}

func (s *CreateUserAddrCommandHandler) Execute(ctx context.Context, req *usermodel.CreateUserAddrReq) error {
	if err := req.Validate(); err != nil {
		return datatype.ErrBadRequest.WithWrap(err).WithDebug(err.Error())
	}

	usrAddr, err := s.userAddrRepo.FindUsrAddrByCityIdAndAddr(ctx, req.CityId, req.Addr)
	if err != nil {
		if !errors.Is(err, usermodel.ErrUserAddrNotFound) {
			return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
		}
	}
	if usrAddr != nil {
		return datatype.ErrConflict.WithWrap(err).WithDebug(usermodel.ErrDuplicated.Error())
	}

	var userAddr usermodel.UserAddress
	copier.Copy(&userAddr, &req)
	userAddr.Id, _ = uuid.NewV7()
	userAddr.Status = sharedModel.StatusActive // Always set Active Status when insert
	now := time.Now().UTC()
	userAddr.CreatedAt = &now
	userAddr.UpdatedAt = &now

	if err := s.userAddrRepo.InsertUserAddress(ctx, userAddr); err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	// set data to response
	req.Id = userAddr.Id

	return nil
}
