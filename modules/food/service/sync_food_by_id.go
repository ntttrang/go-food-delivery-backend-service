package service

import (
	"context"
	"log"

	"github.com/google/uuid"

	rpcrclient "github.com/ntttrang/go-food-delivery-backend-service/modules/food/infras/repository/rpc-client"
	foodmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/food/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

type ISyncFoodByIdRepo interface {
	IndexFood(ctx context.Context, food *foodmodel.FoodDto) error
	DeleteFood(ctx context.Context, id uuid.UUID) error
}

type IFoodByIdRepo interface {
	FindById(ctx context.Context, id uuid.UUID) (foodmodel.Food, error)
	FindByIds(ctx context.Context, ids []uuid.UUID) ([]foodmodel.FoodInfoDto, error)
}

type IRPCRestaurantRepoForSyncFood interface {
	FindByIds(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]rpcrclient.RPCGetByIdsResponseDTO, error)
}

type IRPCCategoryRepoForSyncFood interface {
	FindByIds(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]rpcrclient.CategoryDto, error)
}

type SyncFoodByIdCommandHandler struct {
	foodRepo          IFoodByIdRepo
	indexRepo         ISyncFoodByIdRepo
	rpcrestaurantRepo IRPCRestaurantRepoForSyncFood
	rpccategoryRepo   IRPCCategoryRepoForSyncFood
}

func NewSyncFoodByIdCommandHandler(foodRepo IFoodByIdRepo, indexRepo ISyncFoodByIdRepo, rpcrestaurantRepo IRPCRestaurantRepoForSyncFood, rpccategoryRepo IRPCCategoryRepoForSyncFood) *SyncFoodByIdCommandHandler {
	return &SyncFoodByIdCommandHandler{
		foodRepo:          foodRepo,
		indexRepo:         indexRepo,
		rpcrestaurantRepo: rpcrestaurantRepo,
		rpccategoryRepo:   rpccategoryRepo,
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

	// Get food ratings
	foodRatings, err := s.foodRepo.FindByIds(ctx, []uuid.UUID{food.Id})
	if err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}
	var foodRatingMap = make(map[uuid.UUID]foodmodel.FoodInfoDto, len(foodRatings))
	for _, fr := range foodRatings {
		foodRatingMap[fr.Id] = fr
	}

	if food.RestaurantId == uuid.Nil {
		return datatype.ErrInternalServerError.WithWrap(foodmodel.ErrRestaurantIdEmpty).WithDebug(foodmodel.ErrRestaurantIdEmpty.Error())
	}
	// Get restaurant info
	restaurantMap, err := s.rpcrestaurantRepo.FindByIds(ctx, []uuid.UUID{food.RestaurantId})
	if err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}
	// Get category info
	categoryMap, err := s.rpccategoryRepo.FindByIds(ctx, []uuid.UUID{food.CategoryId})
	if err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	var foodDtos = foodmodel.FoodDto{
		Id:               food.Id,
		RestaurantId:     food.RestaurantId,
		RestaurantName:   restaurantMap[food.RestaurantId].Name,
		RestaurantLat:    restaurantMap[food.RestaurantId].Lat,
		RestaurantLng:    restaurantMap[food.RestaurantId].Lng,
		ShippingFeePerKm: restaurantMap[food.RestaurantId].ShippingFeePerKm,
		CategoryId:       food.CategoryId,
		CategoryName:     categoryMap[food.CategoryId].Name,
		Name:             food.Name,
		Description:      food.Description,
		Price:            food.Price,
		Images:           food.Images,
		AvgPoint:         foodRatingMap[food.Id].AvgPoint,
		CommentQty:       foodRatingMap[food.Id].CommentQty,
		Status:           food.Status,
		CreatedAt:        food.CreatedAt,
		UpdatedAt:        food.UpdatedAt,
	}

	log.Printf("Reindexing %v foods", foodDtos)

	if err := s.indexRepo.IndexFood(ctx, &foodDtos); err != nil {
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
