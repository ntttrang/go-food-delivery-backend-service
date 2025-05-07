package service

import (
	"context"

	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type IRestaurantSearchRepo interface {
	SearchRestaurants(ctx context.Context, req restaurantmodel.RestaurantSearchReq) (*restaurantmodel.RestaurantSearchRes, error)
	GetRestaurantById(ctx context.Context, id string) (*restaurantmodel.RestaurantSearchResDto, error)
}

type SearchRestaurantQueryHandler struct {
	repo IRestaurantSearchRepo
}

func NewSearchRestaurantQueryHandler(repo IRestaurantSearchRepo) *SearchRestaurantQueryHandler {
	return &SearchRestaurantQueryHandler{
		repo: repo,
	}
}

func (s *SearchRestaurantQueryHandler) Execute(ctx context.Context, req restaurantmodel.RestaurantSearchReq) (*restaurantmodel.RestaurantSearchRes, error) {
	// Check if search repository is available
	if s.repo == nil {
		return &restaurantmodel.RestaurantSearchRes{
			Items:      []restaurantmodel.RestaurantSearchResDto{},
			Pagination: req.PagingDto,
		}, nil
	}

	// Execute search
	result, err := s.repo.SearchRestaurants(ctx, req)
	if err != nil {
		return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	// If no items were found, return an empty array instead of nil
	if result.Items == nil {
		result.Items = make([]restaurantmodel.RestaurantSearchResDto, 0)
	}

	// Ensure pagination is properly set
	result.Pagination = sharedModel.PagingDto{
		Page:  req.Page,
		Limit: req.Limit,
		Total: result.Pagination.Total,
	}

	return result, nil
}
