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
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	Price           float64  `json:"price"`
	CategoryID      string   `json:"category_id,omitempty"`
	RestaurantID    string   `json:"restaurant_id"`
	Status          string   `json:"status"`
	Images          string   `json:"images,omitempty"`
	AvgRating       float64  `json:"avg_rating"`
	RatingCount     int      `json:"rating_count"`
	PopularityScore float64  `json:"popularity_score"`
	DeliveryTime    int      `json:"delivery_time"`
	FreeDelivery    bool     `json:"free_delivery"`
	Cuisines        []string `json:"cuisines,omitempty"`
	CreatedAt       string   `json:"created_at"`
	UpdatedAt       string   `json:"updated_at"`
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

	// Default values for new fields
	// In a real implementation, these would be populated from other sources
	cuisines := []string{}
	if f.CategoryId != uuid.Nil {
		// In a real implementation, we would fetch the category name
		// and add it to the cuisines array
		cuisines = append(cuisines, f.CategoryId.String())
	}

	return FoodDocument{
		ID:              f.Id.String(),
		Name:            f.Name,
		Description:     f.Description,
		Price:           f.Price,
		CategoryID:      f.CategoryId.String(),
		RestaurantID:    f.RestaurantId.String(),
		Status:          f.Status,
		Images:          f.Images,
		AvgRating:       0,     // This would be populated from ratings
		RatingCount:     0,     // This would be populated from ratings
		PopularityScore: 0,     // This would be calculated based on views, orders, etc.
		DeliveryTime:    30,    // Default delivery time in minutes
		FreeDelivery:    false, // Default to no free delivery
		Cuisines:        cuisines,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
	}
}

// FromFoodDocument converts a FoodDocument to a FoodSearchResDto
func FromFoodDocument(doc map[string]interface{}) FoodSearchResDto {
	id, _ := uuid.Parse(doc["id"].(string))
	restaurantId, _ := uuid.Parse(doc["restaurant_id"].(string))

	var categoryId uuid.UUID
	if catId, ok := doc["category_id"].(string); ok && catId != "" {
		categoryId, _ = uuid.Parse(catId)
	}

	// Extract cuisines
	var cuisines []string
	if cuisinesVal, ok := doc["cuisines"].([]interface{}); ok {
		for _, c := range cuisinesVal {
			if cStr, ok := c.(string); ok {
				cuisines = append(cuisines, cStr)
			}
		}
	}

	// Get images
	images := ""
	if imagesVal, ok := doc["images"].(string); ok {
		images = imagesVal
	}

	// Get ratings
	avgRating := 0.0
	if val, ok := doc["avg_rating"].(float64); ok {
		avgRating = val
	}

	ratingCount := 0
	if val, ok := doc["rating_count"].(float64); ok {
		ratingCount = int(val)
	}

	// Get popularity score
	popularityScore := 0.0
	if val, ok := doc["popularity_score"].(float64); ok {
		popularityScore = val
	}

	// Get delivery time
	deliveryTime := 30
	if val, ok := doc["delivery_time"].(float64); ok {
		deliveryTime = int(val)
	}

	// Get free delivery
	freeDelivery := false
	if val, ok := doc["free_delivery"].(bool); ok {
		freeDelivery = val
	}

	return FoodSearchResDto{
		Id:              id,
		Name:            doc["name"].(string),
		Description:     doc["description"].(string),
		Status:          doc["status"].(string),
		Price:           doc["price"].(float64),
		RestaurantId:    restaurantId,
		CategoryId:      categoryId,
		AvgRating:       avgRating,
		RatingCount:     ratingCount,
		Images:          images,
		Cuisines:        cuisines,
		PopularityScore: popularityScore,
		DeliveryTime:    deliveryTime,
		FreeDelivery:    freeDelivery,
	}
}

// Enhanced FoodSearchResDto with additional fields
type FoodSearchResDto struct {
	Id              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Status          string    `json:"status"`
	Price           float64   `json:"price"`
	RestaurantId    uuid.UUID `json:"restaurantId"`
	CategoryId      uuid.UUID `json:"categoryId,omitempty"`
	AvgRating       float64   `json:"avgRating"`
	RatingCount     int       `json:"ratingCount"`
	Images          string    `json:"images,omitempty"`
	Cuisines        []string  `json:"cuisines,omitempty"`
	PopularityScore float64   `json:"popularityScore"`
	DeliveryTime    int       `json:"deliveryTime"`
	FreeDelivery    bool      `json:"freeDelivery"`
	Distance        *float64  `json:"distance,omitempty"` // Distance from user in km
	sharedModel.DateDto
}
