package service

import (
	"context"
	"log"

	"github.com/google/uuid"
	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

// IRestaurantRepo defines the interface for restaurant repository operations
type IRestaurantRepo interface {
	FindById(ctx context.Context, id uuid.UUID) (*restaurantmodel.Restaurant, error)
	FindAll(ctx context.Context) ([]restaurantmodel.Restaurant, error)
}

// IRestaurantIndexRepo defines the interface for restaurant indexing operations
type IRestaurantBulkIndexRepo interface {
	IndexRestaurant(ctx context.Context, restaurant *restaurantmodel.Restaurant) error
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
	// Check if repositories are available
	if s.indexRepo == nil || s.restaurantRepo == nil {
		return datatype.ErrInternalServerError.WithDebug("Elasticsearch functionality is not available. Elasticsearch is not configured.")
	}

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

// SyncRestaurant synchronizes a single restaurant with Elasticsearch
func (s *SyncRestaurantIndexCommandHandler) SyncRestaurant(ctx context.Context, id uuid.UUID) error {
	// Check if repositories are available
	if s.indexRepo == nil || s.restaurantRepo == nil {
		return datatype.ErrInternalServerError.WithDebug("Elasticsearch functionality is not available. Elasticsearch is not configured.")
	}

	restaurant, err := s.restaurantRepo.FindById(ctx, id)
	if err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	if err := s.indexRepo.IndexRestaurant(ctx, restaurant); err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	return nil
}
