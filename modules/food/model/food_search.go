package foodmodel

import (
	"github.com/google/uuid"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

// FoodSearchQuery represents the advanced search query for foods
type FoodSearchQuery struct {
	// Basic search fields
	Keyword     string `json:"keyword" form:"keyword"`         // Search in name and description
	Name        string `json:"name" form:"name"`               // Exact name match
	Description string `json:"description" form:"description"` // Description search

	// Advanced filters
	CategoryIds []string `json:"categoryIds" form:"categoryIds"` // Filter by category IDs
	PriceMin    *float64 `json:"priceMin" form:"priceMin"`       // Minimum price
	PriceMax    *float64 `json:"priceMax" form:"priceMax"`       // Maximum price
	Rating      *float64 `json:"rating" form:"rating"`           // Minimum rating

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
	Items      []FoodSearchResDto    `json:"items"`
	Facets     FoodSearchFacets      `json:"facets"`
	Pagination sharedModel.PagingDto `json:"pagination"`
}

// FoodSearchFacets represents the facets in the search response
type FoodSearchFacets struct {
	Categories  []CategoryFacet   `json:"categories"`
	PriceRanges []PriceRangeFacet `json:"priceRanges"`
	Ratings     []RatingFacet     `json:"ratings"`
}

// CategoryFacet represents a category facet
type CategoryFacet struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// PriceRangeFacet represents a price range facet
type PriceRangeFacet struct {
	Range string `json:"range"`
	Count int    `json:"count"`
}

// RatingFacet represents a rating facet
type RatingFacet struct {
	Rating float64 `json:"rating"`
	Count  int     `json:"count"`
}

// FoodDocument represents a food document in Elasticsearch
type FoodDocument struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	CategoryID  string  `json:"category_id,omitempty"`
	Status      string  `json:"status"`
	Images      string  `json:"images,omitempty"`
	AvgRating   float64 `json:"avg_rating"`
	RatingCount int     `json:"rating_count"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// ToFoodDocument converts a Food to a FoodDocument
func (f *Food) ToFoodDocument() FoodDocument {
	var createdAt, updatedAt string
	if f.CreatedAt != nil {
		createdAt = f.CreatedAt.Format("2006-01-02T15:04:05Z")
	}
	if f.UpdatedAt != nil {
		updatedAt = f.UpdatedAt.Format("2006-01-02T15:04:05Z")
	}

	return FoodDocument{
		ID:          f.Id.String(),
		Name:        f.Name,
		Description: f.Description,
		Price:       f.Price,
		CategoryID:  f.CategoryId.String(),
		Status:      f.Status,
		Images:      f.Images,
		AvgRating:   0, // This would be populated from ratings
		RatingCount: 0, // This would be populated from ratings
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

// FromFoodDocument converts a FoodDocument to a FoodSearchResDto
func FromFoodDocument(doc map[string]interface{}) FoodSearchResDto {
	id, _ := uuid.Parse(doc["id"].(string))

	return FoodSearchResDto{
		Id:          id,
		Name:        doc["name"].(string),
		Description: doc["description"].(string),
		Status:      doc["status"].(string),
		Price:       doc["price"].(float64),
		AvgRating:   doc["avg_rating"].(float64),
		//Images:      doc["images"].(string),
	}
}

// Enhanced FoodSearchResDto with additional fields
type FoodSearchResDto struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Price       float64   `json:"price"`
	AvgRating   float64   `json:"avgRating"`
	Images      string    `json:"images"`
	sharedModel.DateDto
}
