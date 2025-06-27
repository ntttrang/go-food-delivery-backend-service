package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/ntttrang/go-food-delivery-backend-service/modules/cart/infras/repository/rpcclient"
	cartmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/cart/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

// Define DTOs & validate
type CartItemSearchDto struct {
	UserID       *string `json:"userId" form:"userId"`
	RestaurantId *string `json:"restaurantId" form:"restaurantId"`
}

type CartItemListReq struct {
	CartItemSearchDto
}

type CartItemListRes struct {
	Items []CartItemDto `json:"items"`
}

type CartItemDto struct {
	ID             uuid.UUID           `json:"id"`
	UserID         uuid.UUID           `json:"userId"`
	RestaurantId   uuid.UUID           `json:"restaurantId"`
	RestaurantName string              `json:"restaurantName"`
	FoodID         uuid.UUID           `json:"foodId"`
	FoodName       string              `json:"foodName"`
	Description    string              `json:"description"`
	Images         string              `json:"images"`
	Price          float64             `json:"price"`
	Quantity       int                 `json:"quantity"`
	Status         datatype.CartStatus `json:"status"`
	sharedModel.DateDto
}

// Initialize service
type IListCartItemRepo interface {
	FindByUserIdAndRestaurantId(ctx context.Context, userId, restaurantId uuid.UUID) ([]cartmodel.Cart, error)
}

type IRpcFoodListCartItemRepo interface {
	FindByIds(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]cartmodel.Food, error)
}

type IRpcRestaurantListCartItemRepo interface {
	FindByIds(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]rpcclient.RPCGetByIdsResponseDTO, error)
}

type ListCartItemQueryHandler struct {
	cartRepo          IListCartItemRepo
	rpcFoodRepo       IRpcFoodListCartItemRepo
	rpcRestaurantRepo IRpcRestaurantListCartItemRepo
}

func NewListCartItemQueryHandler(cartRepo IListCartItemRepo, rpcFoodRepo IRpcFoodListCartItemRepo, rpcRestaurantRepo IRpcRestaurantListCartItemRepo) *ListCartItemQueryHandler {
	return &ListCartItemQueryHandler{
		cartRepo:          cartRepo,
		rpcFoodRepo:       rpcFoodRepo,
		rpcRestaurantRepo: rpcRestaurantRepo,
	}
}

// Implement
func (hdl *ListCartItemQueryHandler) Execute(ctx context.Context, req CartItemListReq) (*CartItemListRes, error) {
	carts, err := hdl.cartRepo.FindByUserIdAndRestaurantId(ctx, uuid.MustParse(*req.UserID), uuid.MustParse(*req.RestaurantId))

	if err != nil {
		return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	var items []CartItemDto
	var foodIds []uuid.UUID
	var restaurantIds []uuid.UUID
	for _, cart := range carts {
		items = append(items, CartItemDto{
			ID:           cart.ID,
			UserID:       cart.UserID,
			RestaurantId: cart.RestaurantId,
			FoodID:       cart.FoodID,
			Quantity:     cart.Quantity,
			Status:       cart.Status,
			DateDto:      cart.DateDto,
		})
		foodIds = append(foodIds, cart.FoodID)
		restaurantIds = append(restaurantIds, cart.RestaurantId)
	}

	foodMap, err := hdl.rpcFoodRepo.FindByIds(ctx, foodIds)
	if err != nil {
		return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}
	for i := 0; i < len(items); i++ {
		items[i].FoodName = foodMap[items[i].FoodID].Name
		items[i].Description = foodMap[items[i].FoodID].Description
		items[i].Images = foodMap[items[i].FoodID].Images
		items[i].Price = foodMap[items[i].FoodID].Price
	}

	restaurantMap, err := hdl.rpcRestaurantRepo.FindByIds(ctx, restaurantIds)
	if err != nil {
		return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}
	for i := 0; i < len(items); i++ {
		items[i].RestaurantName = restaurantMap[items[i].RestaurantId].Name
	}

	var resp CartItemListRes
	resp.Items = items
	return &resp, nil
}
