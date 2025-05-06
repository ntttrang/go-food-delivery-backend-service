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
	IndexFood(ctx context.Context, food *foodmodel.Food) error
	DeleteFood(ctx context.Context, id uuid.UUID) error
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

// Initialize initializes the Elasticsearch index
func (s *SyncFoodIndexCommandHandler) Initialize(ctx context.Context) error {
	// Check if repository is available
	if s.indexRepo == nil {
		return datatype.ErrInternalServerError.WithDebug("Elasticsearch functionality is not available. Elasticsearch is not configured.")
	}

	if err := s.indexRepo.Initialize(ctx); err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}
	return nil
}

// SyncFood synchronizes a single food item with Elasticsearch
func (s *SyncFoodIndexCommandHandler) SyncFood(ctx context.Context, id uuid.UUID) error {
	// Check if repositories are available
	if s.indexRepo == nil || s.foodRepo == nil {
		return datatype.ErrInternalServerError.WithDebug("Elasticsearch functionality is not available. Elasticsearch is not configured.")
	}

	food, err := s.foodRepo.FindById(ctx, id)
	if err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	if err := s.indexRepo.IndexFood(ctx, &food); err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	return nil
}

// DeleteFood removes a food item from Elasticsearch
func (s *SyncFoodIndexCommandHandler) DeleteFood(ctx context.Context, id uuid.UUID) error {
	// Check if repository is available
	if s.indexRepo == nil {
		return datatype.ErrInternalServerError.WithDebug("Elasticsearch functionality is not available. Elasticsearch is not configured.")
	}

	if err := s.indexRepo.DeleteFood(ctx, id); err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}
	return nil
}

// ReindexAll reindexes all foods from the database
func (s *SyncFoodIndexCommandHandler) ReindexAll(ctx context.Context) error {
	// Check if repositories are available
	if s.indexRepo == nil || s.foodRepo == nil {
		return datatype.ErrInternalServerError.WithDebug("Reindex functionality is not available. Elasticsearch is not configured.")
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
