package service

import (
	"context"
	"log"

	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

// IRestaurantRepo defines the interface for restaurant repository operations
type IRestaurantRepo interface {
	FindAll(ctx context.Context) ([]restaurantmodel.Restaurant, error)
}

// IRestaurantIndexRepo defines the interface for restaurant indexing operations
type IRestaurantBulkIndexRepo interface {
	ReindexAllRestaurants(ctx context.Context, restaurants []restaurantmodel.Restaurant) error
}

// SyncRestaurantIndexCommandHandler handles commands to sync restaurant data with Elasticsearch
type SyncRestaurantIndexCommandHandler struct {
	restaurantRepo IRestaurantRepo
	indexRepo      IRestaurantBulkIndexRepo
}

// NewSyncRestaurantIndexCommandHandler creates a new SyncRestaurantIndexCommandHandler
func NewSyncRestaurantIndexCommandHandler(
	restaurantRepo IRestaurantRepo,
	indexRepo IRestaurantBulkIndexRepo,
) *SyncRestaurantIndexCommandHandler {
	return &SyncRestaurantIndexCommandHandler{
		restaurantRepo: restaurantRepo,
		indexRepo:      indexRepo,
	}
}

// SyncAll synchronizes all restaurants with Elasticsearch
func (s *SyncRestaurantIndexCommandHandler) SyncAll(ctx context.Context) error {
	// Get all restaurants from the database
	restaurants, err := s.restaurantRepo.FindAll(ctx)
	if err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	log.Printf("Syncing %d restaurants to Elasticsearch", len(restaurants))

	// Reindex all restaurants
	if err := s.indexRepo.ReindexAllRestaurants(ctx, restaurants); err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	return nil
}
