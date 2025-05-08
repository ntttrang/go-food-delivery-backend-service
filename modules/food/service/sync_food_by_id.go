package service

import (
	"context"

	"github.com/google/uuid"
	foodmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/food/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

type ISyncFoodByIdRepo interface {
	IndexFood(ctx context.Context, food *foodmodel.Food) error
	DeleteFood(ctx context.Context, id uuid.UUID) error
}

type IFoodByIdRepo interface {
	FindById(ctx context.Context, id uuid.UUID) (foodmodel.Food, error)
}

type SyncFoodByIdCommandHandler struct {
	foodRepo  IFoodByIdRepo
	indexRepo ISyncFoodByIdRepo
}

func NewSyncFoodByIdCommandHandler(foodRepo IFoodByIdRepo, indexRepo ISyncFoodByIdRepo) *SyncFoodByIdCommandHandler {
	return &SyncFoodByIdCommandHandler{
		foodRepo:  foodRepo,
		indexRepo: indexRepo,
	}
}

// SyncFood synchronizes a single food item with Elasticsearch
func (s *SyncFoodByIdCommandHandler) SyncFood(ctx context.Context, id uuid.UUID) error {
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
func (s *SyncFoodByIdCommandHandler) DeleteFood(ctx context.Context, id uuid.UUID) error {
	// Check if repository is available
	if s.indexRepo == nil {
		return datatype.ErrInternalServerError.WithDebug("Elasticsearch functionality is not available. Elasticsearch is not configured.")
	}

	if err := s.indexRepo.DeleteFood(ctx, id); err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}
	return nil
}
