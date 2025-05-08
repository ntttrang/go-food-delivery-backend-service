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

// Define DTOs & validate
type UserDetailReq struct {
	Id uuid.UUID `json:"id"`
}

type UserSearchResDto struct {
	Id        uuid.UUID          `json:"id"`
	FirstName string             `json:"firstName"`
	LastName  string             `json:"lastName"`
	Role      usermodel.UserRole `json:"role"`
	Email     string             `json:"email"`
	Phone     string             `json:"phone"`
	CreatedAt *time.Time         `json:"createdAt"`
	UpdatedAt *time.Time         `json:"updatedAt"`
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

	if user.Status == sharedModel.StatusDelete {
		return UserSearchResDto{}, datatype.ErrDeleted.WithError(usermodel.ErrUserDeletedOrBanned.Error())
	}

	var resp UserSearchResDto
	if err := copier.Copy(&resp, &user); err != nil {
		return UserSearchResDto{}, datatype.ErrInternalServerError.WithWrap(errors.New("copier libraries failed"))
	}
	return resp, nil
}
