package service

import (
	"context"

	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/internal/model"
)

type IRestaurantRepository interface {
	Insert(ctx context.Context, data restaurantmodel.Restaurant) error
}

type IRestaurantFoodRepository interface {
	BulkInsert(ctx context.Context, datas []restaurantmodel.RestaurantFood) error
}

type RestaurantService struct {
	resRepo     IRestaurantRepository
	resFoodRepo IRestaurantFoodRepository
}

func NewRestaurantService(resRepo IRestaurantRepository, resFoodRepo IRestaurantFoodRepository) *RestaurantService {
	return &RestaurantService{
		resRepo:     resRepo,
		resFoodRepo: resFoodRepo,
	}
}
