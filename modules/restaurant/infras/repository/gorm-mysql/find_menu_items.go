package restaurantgormmysql

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

func (r *RestaurantFoodRepo) ListMenuItem(ctx context.Context, restaurantId uuid.UUID) ([]restaurantmodel.MenuItemListDto, error) {
	var menuItems []restaurantmodel.MenuItemListDto

	restaurantFoods, err := r.FindByRestaurantId(ctx, restaurantId)
	if err != nil {
		return nil, err
	}

	var foodIds []uuid.UUID
	for _, f := range restaurantFoods {
		if f.Status == sharedModel.StatusActive {
			foodIds = append(foodIds, f.FoodId)
			mi := restaurantmodel.MenuItemListDto{
				FoodId:       f.FoodId,
				RestaurantId: f.RestaurantId,
				CreatedAt:    f.CreatedAt,
				UpdatedAt:    f.UpdatedAt,
			}
			menuItems = append(menuItems, mi)
		}
	}

	foodMap, err := r.foodRPCClient.FindByIds(ctx, foodIds)
	if err != nil {
		return nil, err
	}

	fmt.Println(foodMap)

	return menuItems, nil
}
