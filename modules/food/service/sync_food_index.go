package service

import (
	"context"
	"log"

	"github.com/google/uuid"
	foodmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/food/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

type IFoodIndexRepo interface {
	Initialize(ctx context.Context) error
	ReindexAllFoods(ctx context.Context, foods []foodmodel.Food) error
}

type IFoodRepo interface {
	FindAll(ctx context.Context) ([]foodmodel.Food, error)
	FindById(ctx context.Context, id uuid.UUID) (foodmodel.Food, error)
}

type SyncFoodIndexCommandHandler struct {
	foodRepo  IFoodRepo
	indexRepo IFoodIndexRepo
}

func NewSyncFoodIndexCommandHandler(foodRepo IFoodRepo, indexRepo IFoodIndexRepo) *SyncFoodIndexCommandHandler {
	return &SyncFoodIndexCommandHandler{
		foodRepo:  foodRepo,
		indexRepo: indexRepo,
	}
}

// SyncAll synchronizes all foods with Elasticsearch
func (s *SyncFoodIndexCommandHandler) SyncAll(ctx context.Context) error {
	// Check if repositories are available
	if s.indexRepo == nil || s.foodRepo == nil {
		return datatype.ErrInternalServerError.WithDebug("Sync functionality is not available. Elasticsearch is not configured.")
	}

	// Initialize the index first
	if err := s.indexRepo.Initialize(ctx); err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	// Get all foods from the database
	foods, err := s.foodRepo.FindAll(ctx)
	if err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	log.Printf("Reindexing %d foods", len(foods))

	// Reindex all foods
	if err := s.indexRepo.ReindexAllFoods(ctx, foods); err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	return nil
}
