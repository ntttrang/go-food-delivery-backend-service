package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	usermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/user/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharedatatype "github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

// Define DTOs & validate
type CreateUserAddrReq struct {
	CityId int      `json:"cityId"`
	Title  *string  `json:"title"`
	Icon   string   `json:"icon"`
	Addr   string   `json:"addr"`
	Lat    *float64 `json:"lat"`
	Lng    *float64 `json:"lng"`

	UserId uuid.UUID `json:"-"`
	Id     uuid.UUID `json:"-"`
}

func (r CreateUserAddrReq) TableName() string {
	return usermodel.UserAddress{}.TableName()
}

func (r *CreateUserAddrReq) Validate() error {
	r.Addr = strings.TrimSpace(r.Addr)
	if r.Addr == "" {
		return usermodel.ErrAddrRequired
	}
	return nil
}

// Initilize service
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

// Implement
func (s *CreateUserAddrCommandHandler) Execute(ctx context.Context, req *CreateUserAddrReq) error {
	if err := req.Validate(); err != nil {
		return datatype.ErrBadRequest.WithWrap(err).WithDebug(err.Error())
	}

	usrAddr, err := s.userAddrRepo.FindUsrAddrByCityIdAndAddr(ctx, req.CityId, req.Addr)
	if err != nil {
		if !errors.Is(err, usermodel.ErrUserAddrNotFound) {
			return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
		}
	}
	if usrAddr != nil && usrAddr.UserId == req.UserId {
		return datatype.ErrConflict.WithWrap(err).WithDebug(usermodel.ErrDuplicated.Error())
	}

	var userAddr = usermodel.UserAddress{
		UserId: req.UserId,
		CityId: req.CityId,
		Title:  req.Title,
		Icon:   req.Icon,
		Addr:   req.Addr,
		Lat:    req.Lat,
		Lng:    req.Lng,
		Status: sharedatatype.RecordStatusActive,
	}
	userAddr.Id, _ = uuid.NewV7()
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
