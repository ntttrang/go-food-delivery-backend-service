package service

import (
	"context"

	"github.com/google/uuid"
	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type IUserRepo interface {
	FindById(ctx context.Context, id uuid.UUID) (*restaurantmodel.User, error)
}

type ICreateRestaurantRepository interface {
	Insert(ctx context.Context, restaurant restaurantmodel.Restaurant, restaurantFoods []restaurantmodel.RestaurantFood) error
}

type IBulkCreateRestaurantFoodRepository interface {
	BulkInsert(ctx context.Context, data []restaurantmodel.RestaurantFood) error
}

type CreateCommandHandler struct {
	createRestaurantRepo   ICreateRestaurantRepository
	bulkRestaurantFoodRepo IBulkCreateRestaurantFoodRepository
}

func NewCreateCommandHandler(createRestaurantRepo ICreateRestaurantRepository, bulkRestaurantFoodRepo IBulkCreateRestaurantFoodRepository) *CreateCommandHandler {
	return &CreateCommandHandler{
		createRestaurantRepo:   createRestaurantRepo,
		bulkRestaurantFoodRepo: bulkRestaurantFoodRepo,
	}
}

// func (s *CreateCommandHandler) Execute(ctx context.Context, req *restaurantmodel.RestaurantInsertDto) error {
// 	if err := req.Validate(); err != nil {
// 		return err
// 	}

// 	restaurant := req.ConvertToRestaurant()
// 	restaurant.Id, _ = uuid.NewV7()
// 	restaurant.Status = sharedModel.StatusActive // Always set Active Status when insert

// 	if err := s.createRestaurantRepo.Insert(ctx, *restaurant); err != nil {
// 		return err
// 	}

// 	foods := req.Foods
// 	if len(foods) > 0 {
// 		var restaurantFoods []restaurantmodel.RestaurantFood
// 		for _, f := range foods {
// 			restaurantFood := f.ConvertToRestaurantFood()
// 			restaurantFood.RestaurantID = restaurant.Id
// 			restaurantFoods = append(restaurantFoods, *restaurantFood)
// 		}
// 		if err := s.bulkRestaurantFoodRepo.BulkInsert(ctx, restaurantFoods); err != nil {
// 			return err
// 		}
// 	}

// 	// set data to response
// 	req.Id = restaurant.Id

// 	return nil
// }

func (s *CreateCommandHandler) Execute(ctx context.Context, req *restaurantmodel.RestaurantInsertDto) error {
	if err := req.Validate(); err != nil {
		return err
	}

	// Validate data

	restaurant := req.ConvertToRestaurant()
	restaurant.Id, _ = uuid.NewV7()
	restaurant.Status = sharedModel.StatusActive // Always set Active Status when insert

	foods := req.Foods
	var restaurantFoods []restaurantmodel.RestaurantFood
	if len(foods) > 0 {
		for _, f := range foods {
			restaurantFood := f.ConvertToRestaurantFood()
			restaurantFood.RestaurantID = restaurant.Id
			restaurantFoods = append(restaurantFoods, *restaurantFood)
		}
	}

	if err := s.createRestaurantRepo.Insert(ctx, *restaurant, restaurantFoods); err != nil {
		return err
	}

	// set data to response
	req.Id = restaurant.Id

	return nil
}
