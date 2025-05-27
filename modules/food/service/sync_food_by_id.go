package service

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
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

	var foodDtos foodmodel.FoodDto
	if err := copier.Copy(&foodDtos, &food); err != nil {
		return datatype.ErrInternalServerError.WithWrap(errors.New("copier libraries failed"))
	}

	foodDtos.RestaurantName = restaurantMap[foodDtos.RestaurantId].Name
	foodDtos.RestaurantLat = restaurantMap[foodDtos.RestaurantId].Lat
	foodDtos.RestaurantLng = restaurantMap[foodDtos.RestaurantId].Lng
	foodDtos.ShippingFeePerKm = restaurantMap[foodDtos.RestaurantId].ShippingFeePerKm

	foodDtos.CategoryName = categoryMap[foodDtos.CategoryId].Name

	foodDtos.AvgPoint = foodRatingMap[foodDtos.Id].AvgPoint
	foodDtos.CommentQty = foodRatingMap[foodDtos.Id].CommentQty

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
