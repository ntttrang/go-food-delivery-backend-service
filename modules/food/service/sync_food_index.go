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

type IFoodIndexRepo interface {
	Initialize(ctx context.Context) error
	ReindexAllFoods(ctx context.Context, foods []foodmodel.FoodDto) error
}

type IFoodRepo interface {
	FindAll(ctx context.Context) ([]foodmodel.Food, error)
	FindByIds(ctx context.Context, ids []uuid.UUID) ([]foodmodel.FoodInfoDto, error)
}

type IRPCRestaurantRepo interface {
	FindByIds(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]rpcrclient.RPCGetByIdsResponseDTO, error)
}

type IRPCCategoryRepo interface {
	FindByIds(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]rpcrclient.CategoryDto, error)
}

type SyncFoodIndexCommandHandler struct {
	foodRepo          IFoodRepo
	indexRepo         IFoodIndexRepo
	rpcrestaurantRepo IRPCRestaurantRepo
	rpccategoryRepo   IRPCCategoryRepo
}

func NewSyncFoodIndexCommandHandler(foodRepo IFoodRepo, indexRepo IFoodIndexRepo, rpcrestaurantRepo IRPCRestaurantRepo, rpccategoryRepo IRPCCategoryRepo) *SyncFoodIndexCommandHandler {
	return &SyncFoodIndexCommandHandler{
		foodRepo:          foodRepo,
		indexRepo:         indexRepo,
		rpcrestaurantRepo: rpcrestaurantRepo,
		rpccategoryRepo:   rpccategoryRepo,
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

	var categoryIds []uuid.UUID
	var restaurantIds []uuid.UUID
	var foodIds []uuid.UUID
	var foodDtos []foodmodel.FoodDto
	for _, f := range foods {
		categoryIds = append(categoryIds, f.CategoryId)
		restaurantIds = append(restaurantIds, f.RestaurantId)
		foodIds = append(foodIds, f.Id)

		var dto foodmodel.FoodDto
		if err := copier.Copy(&dto, &f); err != nil {
			return datatype.ErrInternalServerError.WithWrap(errors.New("copier libraries failed"))
		}
		foodDtos = append(foodDtos, dto)
	}

	// Get food ratings
	foodRatings, err := s.foodRepo.FindByIds(ctx, foodIds)
	if err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}
	var foodRatingMap = make(map[uuid.UUID]foodmodel.FoodInfoDto, len(foodRatings))
	for _, fr := range foodRatings {
		foodRatingMap[fr.Id] = fr
	}

	// Get restaurant info
	restaurantMap, err := s.rpcrestaurantRepo.FindByIds(ctx, restaurantIds)
	if err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}
	// Get category info
	categoryMap, err := s.rpccategoryRepo.FindByIds(ctx, categoryIds)
	if err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	for i := 0; i < len(foodDtos); i++ {
		foodDtos[i].RestaurantName = restaurantMap[foodDtos[i].RestaurantId].Name
		foodDtos[i].RestaurantLat = restaurantMap[foodDtos[i].RestaurantId].Lat
		foodDtos[i].RestaurantLng = restaurantMap[foodDtos[i].RestaurantId].Lng
		foodDtos[i].ShippingFeePerKm = restaurantMap[foodDtos[i].RestaurantId].ShippingFeePerKm

		foodDtos[i].CategoryName = categoryMap[foodDtos[i].CategoryId].Name

		foodDtos[i].AvgPoint = foodRatingMap[foodDtos[i].Id].AvgPoint
		foodDtos[i].CommentQty = foodRatingMap[foodDtos[i].Id].CommentQty

	}

	log.Printf("Reindexing %d foods", len(foodDtos))
	// Reindex all foods
	if err := s.indexRepo.ReindexAllFoods(ctx, foodDtos); err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	return nil
}
