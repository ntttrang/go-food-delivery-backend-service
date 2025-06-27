package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

// Define DTOs & validate
type RestaurantDetailReq struct {
	Id uuid.UUID `json:"id"`
}

type RestaurantDetailRes struct {
	Id               uuid.UUID  `json:"id"`
	Name             string     `json:"name"`
	Addr             string     `json:"addr"`
	Logo             string     `json:"logo"`
	ShippingFeePerKm float64    `json:"shippingFeePerKm"`
	Status           string     `json:"status"`
	CreatedAt        *time.Time `json:"createdAt"`
	UpdatedAt        *time.Time `json:"updatedAt"`
}

// Initialize service
type IGetDetailQueryRepo interface {
	FindById(ctx context.Context, id uuid.UUID) (*restaurantmodel.Restaurant, error)
}

type GetDetailQueryHandler struct {
	repo IGetDetailQueryRepo
}

func NewGetDetailQueryHandler(repo IGetDetailQueryRepo) *GetDetailQueryHandler {
	return &GetDetailQueryHandler{repo: repo}
}

// Implement
func (hdl *GetDetailQueryHandler) Execute(ctx context.Context, req RestaurantDetailReq) (*RestaurantDetailRes, error) {
	restaurant, err := hdl.repo.FindById(ctx, req.Id)

	if err != nil {
		if errors.Is(err, restaurantmodel.ErrRestaurantNotFound) {
			return nil, datatype.ErrNotFound.WithDebug(restaurantmodel.ErrRestaurantNotFound.Error())
		}
		return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	if restaurant.Status == string(datatype.StatusDeleted) {
		return nil, datatype.ErrDeleted.WithError(restaurantmodel.ErrRestaurantIsDeleted.Error())
	}

	var resp = RestaurantDetailRes{
		Id:               restaurant.Id,
		Name:             restaurant.Name,
		Addr:             restaurant.Addr,
		Logo:             restaurant.Logo,
		ShippingFeePerKm: restaurant.ShippingFeePerKm,
		Status:           restaurant.Status,
		CreatedAt:        restaurant.CreatedAt,
		UpdatedAt:        restaurant.UpdatedAt,
	}
	return &resp, nil
}
