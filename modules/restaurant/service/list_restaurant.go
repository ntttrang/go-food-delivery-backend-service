package service

import (
	"context"

	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

// Define DTOs & validate
type RestaurantSearchDto struct {
	// CategoryName *string `json:"categoryName"`
	OwnerId *string `json:"ownerId" form:"ownerId"`
	CityId  *int    `json:"cityId" form:"cityId"`
	Status  *string `json:"status" form:"status"`
}

type RestaurantListReq struct {
	RestaurantSearchDto
	sharedModel.PagingDto
	sharedModel.SortingDto
}

// Initialize service
type IListRestaurantRepo interface {
	List(ctx context.Context, req RestaurantListReq) ([]RestaurantSearchResDto, int64, error)
}

type ICategoryRepo interface {
}

type ListQueryHandler struct {
	restaurantRepo IListRestaurantRepo
	categoryRepo   ICategoryRepo
}

func NewListQueryHandler(restaurantRepo IListRestaurantRepo, categoryRepo ICategoryRepo) *ListQueryHandler {
	return &ListQueryHandler{
		restaurantRepo: restaurantRepo,
		categoryRepo:   categoryRepo,
	}
}

// Implement
func (hdl *ListQueryHandler) Execute(ctx context.Context, req RestaurantListReq) (*RestaurantSearchRes, error) {
	restaurants, total, err := hdl.restaurantRepo.List(ctx, req)

	if err != nil {
		return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	var resp RestaurantSearchRes
	resp.Items = restaurants
	resp.Pagination = sharedModel.PagingDto{
		Page:  req.Page,
		Limit: req.Limit,
		Total: total,
	}
	return &resp, nil
}
