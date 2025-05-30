package service

import (
	"context"
	"strings"

	"github.com/google/uuid"
	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

// Define DTOs & validate
type RestaurantInsertDto struct {
	OwnerId          uuid.UUID `json:"-"`
	Name             string    `json:"name"`
	Addr             string    `json:"addr"`
	CityId           int       `json:"cityId"`
	Lat              float64   `json:"lat"`
	Lng              float64   `json:"lng"`
	ShippingFeePerKm float64   `json:"shippingFeePerKm"`

	Id uuid.UUID `json:"-"` // Internal BE
}

func (r RestaurantInsertDto) Validate() error {
	r.Name = strings.TrimSpace(r.Name)

	if r.Name == "" {
		return restaurantmodel.ErrNameRequired
	}

	return nil
}

func (r RestaurantInsertDto) ConvertToRestaurant() *restaurantmodel.Restaurant {
	return &restaurantmodel.Restaurant{
		OwnerId:          r.OwnerId,
		Name:             r.Name,
		Addr:             r.Addr,
		CityId:           r.CityId,
		Lat:              r.Lat,
		Lng:              r.Lng,
		ShippingFeePerKm: r.ShippingFeePerKm,
	}
}

// Initialize service
type IUserRepo interface {
	FindById(ctx context.Context, id uuid.UUID) (*restaurantmodel.User, error)
}

type ICreateRestaurantRepository interface {
	Insert(ctx context.Context, restaurant restaurantmodel.Restaurant) error
}

type IBulkCreateRestaurantFoodRepository interface {
	BulkInsert(ctx context.Context, data []restaurantmodel.RestaurantFood) error
}

type CreateCommandHandler struct {
	createRestaurantRepo   ICreateRestaurantRepository
	bulkRestaurantFoodRepo IBulkCreateRestaurantFoodRepository
}

func NewCreateCommandHandler(createRestaurantRepo ICreateRestaurantRepository, bulkRestaurantFoodRepo IBulkCreateRestaurantFoodRepository) *CreateCommandHandler {
	return &CreateCommandHandler{
		createRestaurantRepo:   createRestaurantRepo,
		bulkRestaurantFoodRepo: bulkRestaurantFoodRepo,
	}
}

// Implement
func (s *CreateCommandHandler) Execute(ctx context.Context, req *RestaurantInsertDto) error {
	if err := req.Validate(); err != nil {
		return datatype.ErrBadRequest.WithWrap(err).WithDebug(err.Error())
	}
	restaurant := req.ConvertToRestaurant()
	restaurant.Id, _ = uuid.NewV7()
	restaurant.Status = sharedModel.StatusActive // Always set Active Status when insert

	if err := s.createRestaurantRepo.Insert(ctx, *restaurant); err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	// set data to response
	req.Id = restaurant.Id

	return nil
}
