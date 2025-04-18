package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type IGetDetailQueryRepo interface {
	FindById(ctx context.Context, id uuid.UUID) (*restaurantmodel.Restaurant, error)
}

type GetDetailQueryHandler struct {
	repo IGetDetailQueryRepo
}

func NewGetDetailQueryHandler(repo IGetDetailQueryRepo) *GetDetailQueryHandler {
	return &GetDetailQueryHandler{repo: repo}
}

func (hdl *GetDetailQueryHandler) Execute(ctx context.Context, req restaurantmodel.RestaurantDetailReq) (restaurantmodel.RestaurantDetailRes, error) {
	restaurant, err := hdl.repo.FindById(ctx, req.Id)

	if err != nil {
		if errors.Is(err, restaurantmodel.ErrRestaurantNotFound) {
			return restaurantmodel.RestaurantDetailRes{}, datatype.ErrNotFound.WithDebug(restaurantmodel.ErrRestaurantNotFound.Error())
		}
		return restaurantmodel.RestaurantDetailRes{}, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	if restaurant.Status == sharedModel.StatusDelete {
		return restaurantmodel.RestaurantDetailRes{}, datatype.ErrDeleted.WithError(restaurantmodel.ErrRestaurantIsDeleted.Error())
	}

	var resp restaurantmodel.RestaurantDetailRes
	if err := copier.Copy(&resp, &restaurant); err != nil {
		return restaurantmodel.RestaurantDetailRes{}, datatype.ErrInternalServerError.WithWrap(errors.New("copier libraries failed"))
	}
	return resp, nil
}
