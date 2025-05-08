package service

import (
	"context"

	"github.com/google/uuid"
	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

type IGetRestaurantByIdRepo interface {
	FindById(ctx context.Context, id uuid.UUID) (*restaurantmodel.Restaurant, error)
}

type ISyncRestaurantByIdRepo interface {
	IndexRestaurant(ctx context.Context, restaurant *restaurantmodel.Restaurant) error
}

type SyncRestaurantByIdCommandHandler struct {
	restaurantRepo IGetRestaurantByIdRepo
	indexRepo      ISyncRestaurantByIdRepo
}

func NewSyncRestaurantByIdCommandHandler(
	restaurantRepo IGetRestaurantByIdRepo,
	indexRepo ISyncRestaurantByIdRepo,
) *SyncRestaurantByIdCommandHandler {
	return &SyncRestaurantByIdCommandHandler{
		restaurantRepo: restaurantRepo,
		indexRepo:      indexRepo,
	}
}

// SyncRestaurant synchronizes a single restaurant with Elasticsearch
func (s *SyncRestaurantByIdCommandHandler) SyncRestaurant(ctx context.Context, id uuid.UUID) error {
	restaurant, err := s.restaurantRepo.FindById(ctx, id)
	if err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	if err := s.indexRepo.IndexRestaurant(ctx, restaurant); err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	return nil
}
