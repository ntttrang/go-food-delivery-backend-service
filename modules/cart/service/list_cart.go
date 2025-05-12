package service

import (
	"context"

	"github.com/google/uuid"
	cartmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/cart/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

// Define DTOs & validate
type CartSearchDto struct {
	UserID *string `json:"userId" form:"userId"`
}

type CartListReq struct {
	CartSearchDto
	sharedModel.PagingDto
	sharedModel.SortingDto
}

type CartListRes struct {
	Items      []CartItemDto         `json:"items"`
	Pagination sharedModel.PagingDto `json:"pagination"`
}

type CartItemDto struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"userId"`
	FoodID      uuid.UUID `json:"foodId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Quantity    int       `json:"quantity"`
	Status      string    `json:"status"`
	sharedModel.DateDto
}

// Initialize service
type IListCartRepo interface {
	ListByUserId(ctx context.Context, req CartListReq) ([]cartmodel.Cart, int64, error)
}

type IRpcFoodRepo interface {
	FindByIds(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]cartmodel.Food, error)
}

type ListQueryHandler struct {
	cartRepo    IListCartRepo
	rpcFoodRepo IRpcFoodRepo
}

func NewListQueryHandler(cartRepo IListCartRepo, rpcFoodRepo IRpcFoodRepo) *ListQueryHandler {
	return &ListQueryHandler{
		cartRepo:    cartRepo,
		rpcFoodRepo: rpcFoodRepo,
	}
}

// Implement
func (hdl *ListQueryHandler) Execute(ctx context.Context, req CartListReq) (*CartListRes, error) {
	carts, total, err := hdl.cartRepo.ListByUserId(ctx, req)

	if err != nil {
		return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	var items []CartItemDto
	var foodIds []uuid.UUID
	for _, cart := range carts {
		items = append(items, CartItemDto{
			ID:       cart.ID,
			UserID:   cart.UserID,
			FoodID:   cart.FoodID,
			Quantity: cart.Quantity,
			Status:   cart.Status,
			DateDto:  cart.DateDto,
		})
		foodIds = append(foodIds, cart.FoodID)
	}

	foodMap, err := hdl.rpcFoodRepo.FindByIds(ctx, foodIds)
	if err != nil {
		return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}
	for i := 0; i < len(items); i++ {
		items[i].Name = foodMap[items[i].FoodID].Name
		items[i].Description = foodMap[items[i].FoodID].Description
	}

	var resp CartListRes
	resp.Items = items
	resp.Pagination = sharedModel.PagingDto{
		Page:  req.Page,
		Limit: req.Limit,
		Total: total,
	}
	return &resp, nil
}
