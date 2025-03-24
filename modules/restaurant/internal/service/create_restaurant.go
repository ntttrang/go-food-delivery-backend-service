package service

import (
	"context"

	"github.com/google/uuid"
	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/internal/model"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

func (s *RestaurantService) CreateRestaurant(ctx context.Context, data *restaurantmodel.RestaurantInsertDto) error {
	if err := data.Validate(); err != nil {
		return err
	}

	restaurant := data.ConvertToRestaurant()
	restaurant.Id, _ = uuid.NewV7()
	restaurant.Status = sharedModel.StatusActive // Always set Active Status when insert

	if err := s.resRepo.Insert(ctx, *restaurant); err != nil {
		return err
	}

	foods := data.Foods
	var restaurantFoods []restaurantmodel.RestaurantFood
	for _, f := range foods {
		restaurantFood := f.ConvertToRestaurantFood()
		restaurantFood.RestaurantID = restaurant.Id
		restaurantFoods = append(restaurantFoods, *restaurantFood)
	}
	if err := s.resFoodRepo.BulkInsert(ctx, restaurantFoods); err != nil {
		return err
	}

	// set data to response
	data.Id = restaurant.Id

	return nil
}
