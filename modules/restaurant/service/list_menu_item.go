package service

import (
	"context"

	"github.com/google/uuid"
	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type IRestaurantFoodRepo interface {
	FindByRestaurantId(ctx context.Context, id uuid.UUID) ([]restaurantmodel.RestaurantFood, error)
}
type IRPCFoodRepo interface {
	FindByIds(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]restaurantmodel.Foods, error)
}

type IRPCCategoryRepo interface {
	FindByIds(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]restaurantmodel.CategoryDto, error)
}

type ListMenuItemQueryHandler struct {
	restaurantFoodRepo IRestaurantFoodRepo
	rpcFoodRepo        IRPCFoodRepo
	rpcCategoryRepo    IRPCCategoryRepo
}

func NewListMenuItemQueryHandler(restaurantFoodRepo IRestaurantFoodRepo, rpcFoodRepo IRPCFoodRepo, rpcCategoryRepo IRPCCategoryRepo) *ListMenuItemQueryHandler {
	return &ListMenuItemQueryHandler{
		restaurantFoodRepo: restaurantFoodRepo,
		rpcFoodRepo:        rpcFoodRepo,
		rpcCategoryRepo:    rpcCategoryRepo,
	}
}

func (hdl *ListMenuItemQueryHandler) Execute(ctx context.Context, restaurantId uuid.UUID) (*restaurantmodel.MenuItemListRes, error) {
	var resp restaurantmodel.MenuItemListRes
	var menuItems []restaurantmodel.MenuItemListDto

	restaurantFoods, err := hdl.restaurantFoodRepo.FindByRestaurantId(ctx, restaurantId)
	if err != nil {
		return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
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

	foodMap, err := hdl.rpcFoodRepo.FindByIds(ctx, foodIds)
	if err != nil {
		return nil, err
	}

	var categoryIds []uuid.UUID
	for _, v := range foodMap {
		categoryIds = append(categoryIds, v.CategoryId)
	}

	categoryMap, err := hdl.rpcCategoryRepo.FindByIds(ctx, categoryIds)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(menuItems); i++ {
		menuItems[i].Name = foodMap[menuItems[i].FoodId].Name
		menuItems[i].Description = foodMap[menuItems[i].FoodId].Description
		menuItems[i].ImageURL = foodMap[menuItems[i].FoodId].Images
		menuItems[i].Price = foodMap[menuItems[i].FoodId].Price
		menuItems[i].Point = foodMap[menuItems[i].FoodId].AvgPoint
		menuItems[i].CommentQty = foodMap[menuItems[i].FoodId].CommentQty
		menuItems[i].CategoryId = foodMap[menuItems[i].FoodId].CategoryId
		menuItems[i].CategoryName = categoryMap[menuItems[i].CategoryId].Name
	}

	resp.Items = menuItems
	return &resp, nil
}
