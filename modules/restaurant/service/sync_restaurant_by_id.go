package service

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"

	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharedmodel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type IGetRestaurantByIdRepo interface {
	FindRestaurantByIds(ctx context.Context, ids []uuid.UUID) ([]FindRestaurantByIdsDto, error)
}

type IGetRestaurantFoodRepo interface {
	FindFoodByRestaurantIds(ctx context.Context, ids []uuid.UUID) ([]restaurantmodel.RestaurantFood, error)
}

type IRpcFoodSyncRestaurantbyIdRepo interface {
	FindByIds(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]restaurantmodel.Foods, error)
}

type ISyncRestaurantByIdRepo interface {
	IndexRestaurant(ctx context.Context, restaurant *restaurantmodel.RestaurantInfoDto) error
}

type SyncRestaurantByIdCommandHandler struct {
	restaurantRepo     IGetRestaurantByIdRepo
	restaurantFoodRepo IGetRestaurantFoodRepo
	indexRepo          ISyncRestaurantByIdRepo
	rpcFood            IRpcFoodSyncRestaurantbyIdRepo
}

func NewSyncRestaurantByIdCommandHandler(
	restaurantRepo IGetRestaurantByIdRepo,
	restaurantFoodRepo IGetRestaurantFoodRepo,
	indexRepo ISyncRestaurantByIdRepo,
	rpcFood IRpcFoodSyncRestaurantbyIdRepo,
) *SyncRestaurantByIdCommandHandler {
	return &SyncRestaurantByIdCommandHandler{
		restaurantRepo:     restaurantRepo,
		restaurantFoodRepo: restaurantFoodRepo,
		indexRepo:          indexRepo,
		rpcFood:            rpcFood,
	}
}

// SyncRestaurant synchronizes a single restaurant with Elasticsearch
func (s *SyncRestaurantByIdCommandHandler) SyncRestaurant(ctx context.Context, id uuid.UUID) error {
	restaurants, err := s.restaurantRepo.FindRestaurantByIds(ctx, []uuid.UUID{id})
	if err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	if len(restaurants) == 0 {
		log.Printf("There is no data to sync")
		return nil
	}

	restaurant := restaurants[0]
	var restaurantDto = restaurantmodel.RestaurantInfoDto{
		Restaurant: restaurantmodel.Restaurant{
			Id:               restaurant.Id,
			OwnerId:          restaurant.OwnerId,
			Name:             restaurant.Name,
			Addr:             restaurant.Addr,
			CityId:           restaurant.CityId,
			Lat:              restaurant.Lat,
			Lng:              restaurant.Lng,
			Cover:            restaurant.Cover,
			Logo:             restaurant.Logo,
			ShippingFeePerKm: restaurant.ShippingFeePerKm,
			Status:           restaurant.Status,
			DateDto: sharedmodel.DateDto{
				CreatedAt: restaurant.CreatedAt,
				UpdatedAt: restaurant.UpdatedAt,
			},
		},
		AvgPoint:   restaurant.AvgPoint,
		CommentQty: restaurant.CommentQty,
		LikesQty:   restaurant.LikesQty,
	}

	restaurantId := restaurant.Id
	// Get restaurant food
	restaurantFoods, err := s.restaurantFoodRepo.FindFoodByRestaurantIds(ctx, []uuid.UUID{restaurantId})
	if err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	if len(restaurantFoods) > 0 {
		rf := restaurantFoods[0]
		// Call RPC to get food info
		foodMap, err := s.rpcFood.FindByIds(ctx, []uuid.UUID{rf.FoodId})
		if err != nil {
			return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
		}

		foodInfo := restaurantmodel.FoodInfo{
			CategoryId: foodMap[rf.FoodId].CategoryId,
			FoodId:     foodMap[rf.FoodId].Id,
			FoodName:   foodMap[rf.FoodId].Name,
		}
		fmt.Println(foodInfo)
		restaurantDto.FoodInfos = append(restaurantDto.FoodInfos, foodInfo)
	}

	log.Printf("Syncing restaurant id = %v to Elasticsearch", restaurantDto.Id)

	if err := s.indexRepo.IndexRestaurant(ctx, &restaurantDto); err != nil {
		return datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	return nil
}
