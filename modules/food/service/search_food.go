package service

import (
	"context"

	foodmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/food/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

// Define DTOs & validate
// FoodSearchQuery represents the advanced search query for foods
type FoodSearchQuery struct {
	// Basic search fields
	Keyword     string `json:"keyword" form:"keyword"`         // Search in name and description
	Name        string `json:"name" form:"name"`               // Exact name match
	Description string `json:"description" form:"description"` // Description search

	// Advanced filters
	CategoryIds  []string `json:"categoryIds" form:"categoryIds"`   // Filter by category IDs (cuisines)
	PriceMin     *float64 `json:"priceMin" form:"priceMin"`         // Minimum price
	PriceMax     *float64 `json:"priceMax" form:"priceMax"`         // Maximum price
	Rating       *float64 `json:"rating" form:"rating"`             // Minimum rating
	FreeDelivery *bool    `json:"freeDelivery" form:"freeDelivery"` // Free delivery option

	// Restaurant filters
	RestaurantId *string `json:"restaurantId" form:"restaurantId"` // Filter by restaurant ID

	// Location filters for delivery time calculation
	Lat         *float64 `json:"lat" form:"lat"`                 // User's latitude
	Lng         *float64 `json:"lng" form:"lng"`                 // User's longitude
	MaxDistance *float64 `json:"maxDistance" form:"maxDistance"` // Maximum distance in km

	// Sorting
	sharedModel.SortingDto

	// Pagination
	sharedModel.PagingDto
}

// FoodSearchReq represents the search request for foods
type FoodSearchReq struct {
	FoodSearchQuery
}

// FoodSearchRes represents the search response for foods
type FoodSearchRes struct {
	Items      []foodmodel.FoodSearchResDto `json:"items"`
	Facets     foodmodel.FoodSearchFacets   `json:"facets"`
	Pagination sharedModel.PagingDto        `json:"pagination"`
}

// Initilize service
type IFoodSearchRepo interface {
	SearchFoods(ctx context.Context, req FoodSearchReq) (*FoodSearchRes, error)
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

// Implement
func (s *SearchFoodQueryHandler) Execute(ctx context.Context, req FoodSearchReq) (*FoodSearchRes, error) {

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
