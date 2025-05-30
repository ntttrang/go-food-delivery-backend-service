package service

import (
	"context"
	"log"

	"github.com/google/uuid"
	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
)

// RestaurantElasticsearchHandler implements IRestaurantEventHandler for Elasticsearch operations
type RestaurantElasticsearchHandler struct {
	indexRepo IRestaurantIndexRepo
}

// IRestaurantIndexRepo defines the interface for restaurant indexing operations
type IRestaurantIndexRepo interface {
	IndexRestaurant(ctx context.Context, restaurant *restaurantmodel.RestaurantInfoDto) error
	DeleteRestaurant(ctx context.Context, id uuid.UUID) error
}

// NewRestaurantElasticsearchHandler creates a new RestaurantElasticsearchHandler
func NewRestaurantElasticsearchHandler(indexRepo IRestaurantIndexRepo) *RestaurantElasticsearchHandler {
	return &RestaurantElasticsearchHandler{
		indexRepo: indexRepo,
	}
}

// OnRestaurantCreated handles the restaurant created event
func (h *RestaurantElasticsearchHandler) OnRestaurantCreated(ctx context.Context, restaurant *restaurantmodel.RestaurantInfoDto) error {
	log.Printf("Indexing new restaurant: %s", restaurant.Id)
	return h.indexRepo.IndexRestaurant(ctx, restaurant)
}

// OnRestaurantUpdated handles the restaurant updated event
func (h *RestaurantElasticsearchHandler) OnRestaurantUpdated(ctx context.Context, restaurant *restaurantmodel.RestaurantInfoDto) error {
	log.Printf("Updating indexed restaurant: %s", restaurant.Id)
	return h.indexRepo.IndexRestaurant(ctx, restaurant)
}

// OnRestaurantDeleted handles the restaurant deleted event
func (h *RestaurantElasticsearchHandler) OnRestaurantDeleted(ctx context.Context, id uuid.UUID) error {
	log.Printf("Deleting restaurant from index: %s", id)
	return h.indexRepo.DeleteRestaurant(ctx, id)
}
