package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

// Define DTOs & validate
type RestaurantDetailReq struct {
	Id uuid.UUID `json:"id"`
}

type RestaurantDetailRes struct {
	Id               uuid.UUID       `json:"id"`
	Name             string          `json:"name"`
	Addr             string          `json:"addr"`
	Logo             json.RawMessage `json:"logo"`
	ShippingFeePerKm float64         `json:"shippingFeePerKm"`
	Status           string          `json:"status"`
	//CategoryName     string          `json:"categoryName"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
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

	if restaurant.Status == sharedModel.StatusDelete {
		return nil, datatype.ErrDeleted.WithError(restaurantmodel.ErrRestaurantIsDeleted.Error())
	}

	var resp RestaurantDetailRes
	if err := copier.Copy(&resp, &restaurant); err != nil {
		return nil, datatype.ErrInternalServerError.WithWrap(errors.New("copier libraries failed"))
	}
	return &resp, nil
}
