package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	usermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/user/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

// Define DTOs & validate
type UserDetailReq struct {
	Id uuid.UUID `json:"id"`
}

type UserSearchResDto struct {
	Id        uuid.UUID         `json:"id"`
	FirstName string            `json:"firstName"`
	LastName  string            `json:"lastName"`
	Role      datatype.UserRole `json:"role"`
	Email     string            `json:"email"`
	Phone     string            `json:"phone"`
	Avatar    string            `json:"avatar"`
	CreatedAt *time.Time        `json:"createdAt"`
	UpdatedAt *time.Time        `json:"updatedAt"`
}

// Initilize service
type IGetDetailQueryRepo interface {
	FindById(ctx context.Context, id uuid.UUID) (*usermodel.User, error)
}

type GetDetailQueryHandler struct {
	repo IGetDetailQueryRepo
}

func NewGetDetailQueryHandler(repo IGetDetailQueryRepo) *GetDetailQueryHandler {
	return &GetDetailQueryHandler{repo: repo}
}

// Implement
func (hdl *GetDetailQueryHandler) Execute(ctx context.Context, req UserDetailReq) (UserSearchResDto, error) {
	user, err := hdl.repo.FindById(ctx, req.Id)

	if err != nil {
		if errors.Is(err, usermodel.ErrUserNotFound) {
			return UserSearchResDto{}, datatype.ErrNotFound.WithDebug(usermodel.ErrUserNotFound.Error())
		}
		return UserSearchResDto{}, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	if user.Status == datatype.StatusDeleted {
		return UserSearchResDto{}, datatype.ErrDeleted.WithError(usermodel.ErrUserDeletedOrBanned.Error())
	}

	var resp = UserSearchResDto{
		Id:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		Email:     user.Email,
		Phone:     user.Phone,
		Avatar:    user.Avatar,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	return resp, nil
}
