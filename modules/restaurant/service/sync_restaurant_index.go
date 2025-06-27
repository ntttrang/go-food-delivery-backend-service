package service

import (
	"context"
	"log"

	"github.com/google/uuid"

	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharedmodel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type FindRestaurantByIdsDto struct {
	restaurantmodel.Restaurant
	AvgPoint   float64 `json:"avgPoint"`
	CommentQty int     `json:"commentQty"`
	LikesQty   int     `json:"likesQty"`
}

// IRestaurantRepo defines the interface for restaurant repository operations
type IRestaurantRepo interface {
	FindRestaurantByIds(ctx context.Context, ids []uuid.UUID) ([]FindRestaurantByIdsDto, error)
}

type IRestaurantFoodRestaurantIndexRepo interface {
	FindFoodByRestaurantIds(ctx context.Context, ids []uuid.UUID) ([]restaurantmodel.RestaurantFood, error)
}

// IRestaurantIndexRepo defines the interface for restaurant indexing operations
type IRestaurantBulkIndexRepo interface {
	ReindexAllRestaurants(ctx context.Context, restaurants []restaurantmodel.RestaurantInfoDto) error
}

type IRpcFoodSyncRestaurantIndexRepo interface {
	FindByIds(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]restaurantmodel.Foods, error)
}

// SyncRestaurantIndexCommandHandler handles commands to sync restaurant data with Elasticsearch
type SyncRestaurantIndexCommandHandler struct {
	restaurantRepo     IRestaurantRepo
	restaurantFoodRepo IRestaurantFoodRestaurantIndexRepo
	indexRepo          IRestaurantBulkIndexRepo
	rpcFood            IRpcFoodSyncRestaurantIndexRepo
}

// NewSyncRestaurantIndexCommandHandler creates a new SyncRestaurantIndexCommandHandler
func NewSyncRestaurantIndexCommandHandler(
	restaurantRepo IRestaurantRepo,
	restaurantFoodRepo IRestaurantFoodRestaurantIndexRepo,
	indexRepo IRestaurantBulkIndexRepo,
	rpcFood IRpcFoodSyncRestaurantIndexRepo,
) *SyncRestaurantIndexCommandHandler {
	return &SyncRestaurantIndexCommandHandler{
		restaurantRepo:     restaurantRepo,
		restaurantFoodRepo: restaurantFoodRepo,
		indexRepo:          indexRepo,
		rpcFood:            rpcFood,
	}
}

// SyncAll synchronizes all restaurants with Elasticsearch
func (s *SyncRestaurantIndexCommandHandler) SyncAll(ctx context.Context) error {
	// Get all restaurants from the database
	restaurants, err := s.restaurantRepo.FindRestaurantByIds(ctx, []uuid.UUID{})
	if err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	var foodIds []uuid.UUID
	var restaurantIds []uuid.UUID
	var restaurantDtos []restaurantmodel.RestaurantInfoDto
	for _, r := range restaurants {
		restaurantIds = append(restaurantIds, r.Id)

		var dto = restaurantmodel.RestaurantInfoDto{
			Restaurant: restaurantmodel.Restaurant{
				Id:               r.Id,
				OwnerId:          r.OwnerId,
				Name:             r.Name,
				Addr:             r.Addr,
				CityId:           r.CityId,
				Lat:              r.Lat,
				Lng:              r.Lng,
				Cover:            r.Cover,
				Logo:             r.Logo,
				ShippingFeePerKm: r.ShippingFeePerKm,
				Status:           r.Status,
				DateDto: sharedmodel.DateDto{
					CreatedAt: r.CreatedAt,
					UpdatedAt: r.UpdatedAt,
				},
			},
			AvgPoint:   r.AvgPoint,
			CommentQty: r.CommentQty,
			LikesQty:   r.LikesQty,
		}
		restaurantDtos = append(restaurantDtos, dto)
	}

	// Get restaurant food
	restaurantFoods, err := s.restaurantFoodRepo.FindFoodByRestaurantIds(ctx, restaurantIds)
	if err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	for _, rf := range restaurantFoods {
		foodIds = append(foodIds, rf.FoodId)
	}

	// Call RPC to get food info
	foodMap, err := s.rpcFood.FindByIds(ctx, foodIds)
	if err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	for i := 0; i < len(restaurantFoods); i++ {
		foodInfo := restaurantmodel.FoodInfo{
			CategoryId: foodMap[restaurantFoods[i].FoodId].CategoryId,
			FoodId:     foodMap[restaurantFoods[i].FoodId].Id,
			FoodName:   foodMap[restaurantFoods[i].FoodId].Name,
		}
		restaurantDtos[i].FoodInfos = append(restaurantDtos[i].FoodInfos, foodInfo)
	}
	log.Printf("Syncing %d restaurants to Elasticsearch", len(restaurantDtos))

	// Reindex all restaurants
	if err := s.indexRepo.ReindexAllRestaurants(ctx, restaurantDtos); err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	return nil
}
