package service

import (
	"context"
	"log"

	"github.com/google/uuid"

	foodmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/food/model"
)

// FoodElasticsearchHandler implements IFoodEventHandler for Elasticsearch operations
type FoodElasticsearchHandler struct {
	indexRepo ISyncFoodByIdRepo
}

// NewFoodElasticsearchHandler creates a new FoodElasticsearchHandler
func NewFoodElasticsearchHandler(indexRepo ISyncFoodByIdRepo) *FoodElasticsearchHandler {
	return &FoodElasticsearchHandler{
		indexRepo: indexRepo,
	}
}

// OnFoodCreated handles the food created event
func (h *FoodElasticsearchHandler) OnFoodCreated(ctx context.Context, food *foodmodel.FoodDto) error {
	log.Printf("Indexing new food: %s", food.Id)
	return h.indexRepo.IndexFood(ctx, food)
}

// OnFoodUpdated handles the food updated event
func (h *FoodElasticsearchHandler) OnFoodUpdated(ctx context.Context, food *foodmodel.FoodDto) error {
	log.Printf("Updating indexed food: %s", food.Id)
	return h.indexRepo.IndexFood(ctx, food)
}

// OnFoodDeleted handles the food deleted event
func (h *FoodElasticsearchHandler) OnFoodDeleted(ctx context.Context, id uuid.UUID) error {
	log.Printf("Deleting food from index: %s", id)
	return h.indexRepo.DeleteFood(ctx, id)
}
