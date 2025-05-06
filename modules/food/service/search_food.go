package service

import (
	"context"

	foodmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/food/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type IFoodSearchRepo interface {
	SearchFoods(ctx context.Context, req foodmodel.FoodSearchReq) (*foodmodel.FoodSearchRes, error)
	GetFoodById(ctx context.Context, id string) (*foodmodel.FoodSearchResDto, error)
}

type SearchFoodQueryHandler struct {
	repo IFoodSearchRepo
}

func NewSearchFoodQueryHandler(repo IFoodSearchRepo) *SearchFoodQueryHandler {
	return &SearchFoodQueryHandler{
		repo: repo,
	}
}

func (s *SearchFoodQueryHandler) Execute(ctx context.Context, req foodmodel.FoodSearchReq) (*foodmodel.FoodSearchRes, error) {

	// Execute search
	result, err := s.repo.SearchFoods(ctx, req)
	if err != nil {
		return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	// If no items were found, return an empty array instead of nil
	if result.Items == nil {
		result.Items = make([]foodmodel.FoodSearchResDto, 0)
	}

	// Ensure pagination is properly set
	result.Pagination = sharedModel.PagingDto{
		Page:  req.Page,
		Limit: req.Limit,
		Total: result.Pagination.Total,
	}

	return result, nil
}
