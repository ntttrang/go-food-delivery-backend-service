package restaurantmodel

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

// RestaurantDocument represents a restaurant document in Elasticsearch
type RestaurantDocument struct {
	ID               string          `json:"id"`
	Name             string          `json:"name"`
	Address          string          `json:"address"`
	CityID           int             `json:"city_id"`
	Lat              float64         `json:"lat"`
	Lng              float64         `json:"lng"`
	Logo             json.RawMessage `json:"logo,omitempty"`
	Cover            json.RawMessage `json:"cover,omitempty"`
	ShippingFeePerKm float64         `json:"shipping_fee_per_km"`
	Status           string          `json:"status"`
	AvgRating        float64         `json:"avg_rating"`
	RatingCount      int             `json:"rating_count"`
	Cuisines         []string        `json:"cuisines,omitempty"` // Categories/cuisines associated with this restaurant
	PopularityScore  float64         `json:"popularity_score"`   // Based on views, orders, etc.
	DeliveryTime     int             `json:"delivery_time"`      // Estimated delivery time in minutes
	CreatedAt        string          `json:"created_at"`
	UpdatedAt        string          `json:"updated_at"`
}

// RestaurantSearchQuery represents the advanced search query for restaurants
type RestaurantSearchQuery struct {
	// Basic search fields
	Keyword string `json:"keyword" form:"keyword"` // Search in name and address
	Name    string `json:"name" form:"name"`       // Exact name match

	// Location filters
	CityID *int     `json:"cityId" form:"cityId"` // Filter by city ID
	Lat    *float64 `json:"lat" form:"lat"`       // User's latitude for distance calculation
	Lng    *float64 `json:"lng" form:"lng"`       // User's longitude for distance calculation
	Radius *float64 `json:"radius" form:"radius"` // Search radius in km

	// Advanced filters
	Cuisines    []string `json:"cuisines" form:"cuisines"`       // Filter by cuisine types
	Rating      *float64 `json:"rating" form:"rating"`           // Minimum rating
	PriceMin    *float64 `json:"priceMin" form:"priceMin"`       // Minimum price
	PriceMax    *float64 `json:"priceMax" form:"priceMax"`       // Maximum price
	FreeShipping *bool    `json:"freeShipping" form:"freeShipping"` // Free shipping option

	// Sorting
	sharedModel.SortingDto

	// Pagination
	sharedModel.PagingDto
}

// RestaurantSearchReq represents the search request for restaurants
type RestaurantSearchReq struct {
	RestaurantSearchQuery
}

// RestaurantSearchRes represents the search response for restaurants
type RestaurantSearchRes struct {
	Items      []RestaurantSearchResDto `json:"items"`
	Facets     RestaurantSearchFacets   `json:"facets"`
	Pagination sharedModel.PagingDto    `json:"pagination"`
}

// RestaurantSearchFacets represents the facets in the search response
type RestaurantSearchFacets struct {
	Cuisines     []CuisineFacet    `json:"cuisines"`
	PriceRanges  []PriceRangeFacet `json:"priceRanges"`
	Ratings      []RatingFacet     `json:"ratings"`
	DeliveryTime []TimeFacet       `json:"deliveryTime"`
}

// CuisineFacet represents a cuisine facet
type CuisineFacet struct {
	Value string `json:"value"`
	Count int    `json:"count"`
}

// PriceRangeFacet represents a price range facet
type PriceRangeFacet struct {
	From  *float64 `json:"from,omitempty"`
	To    *float64 `json:"to,omitempty"`
	Count int      `json:"count"`
}

// RatingFacet represents a rating facet
type RatingFacet struct {
	From  float64 `json:"from"`
	To    float64 `json:"to"`
	Count int     `json:"count"`
}

// TimeFacet represents a delivery time facet
type TimeFacet struct {
	From  int `json:"from"`
	To    int `json:"to"`
	Count int `json:"count"`
}

// RestaurantSearchResDto represents a restaurant in the search results
type RestaurantSearchResDto struct {
	ID               uuid.UUID       `json:"id"`
	Name             string          `json:"name"`
	Address          string          `json:"address"`
	Logo             json.RawMessage `json:"logo,omitempty"`
	Cover            json.RawMessage `json:"cover,omitempty"`
	ShippingFeePerKm float64         `json:"shippingFeePerKm"`
	AvgRating        float64         `json:"avgRating"`
	RatingCount      int             `json:"ratingCount"`
	Cuisines         []string        `json:"cuisines,omitempty"`
	PopularityScore  float64         `json:"popularityScore"`
	DeliveryTime     int             `json:"deliveryTime"`
	Distance         *float64        `json:"distance,omitempty"` // Distance from user in km
	Status           string          `json:"status"`
	CreatedAt        *time.Time      `json:"createdAt,omitempty"`
	UpdatedAt        *time.Time      `json:"updatedAt,omitempty"`
}

// ToRestaurantDocument converts a Restaurant to a RestaurantDocument for Elasticsearch
func (r *Restaurant) ToRestaurantDocument() RestaurantDocument {
	createdAt := ""
	updatedAt := ""
	
	if r.CreatedAt != nil {
		createdAt = r.CreatedAt.Format(time.RFC3339)
	}
	
	if r.UpdatedAt != nil {
		updatedAt = r.UpdatedAt.Format(time.RFC3339)
	}
	
	return RestaurantDocument{
		ID:               r.Id.String(),
		Name:             r.Name,
		Address:          r.Addr,
		CityID:           r.CityId,
		Lat:              r.Lat,
		Lng:              r.Lng,
		Logo:             r.Logo,
		Cover:            r.Cover,
		ShippingFeePerKm: r.ShippingFeePerKm,
		Status:           r.Status,
		CreatedAt:        createdAt,
		UpdatedAt:        updatedAt,
		// These fields would be populated from other sources
		AvgRating:       0,
		RatingCount:     0,
		PopularityScore: 0,
		DeliveryTime:    30, // Default delivery time
	}
}

// FromRestaurantDocument converts an Elasticsearch document to a RestaurantSearchResDto
func FromRestaurantDocument(doc map[string]interface{}) RestaurantSearchResDto {
	id, _ := uuid.Parse(doc["id"].(string))
	
	var logo, cover json.RawMessage
	if logoVal, ok := doc["logo"]; ok && logoVal != nil {
		logoBytes, _ := json.Marshal(logoVal)
		logo = logoBytes
	}
	
	if coverVal, ok := doc["cover"]; ok && coverVal != nil {
		coverBytes, _ := json.Marshal(coverVal)
		cover = coverBytes
	}
	
	// Parse cuisines
	var cuisines []string
	if cuisinesVal, ok := doc["cuisines"]; ok && cuisinesVal != nil {
		for _, c := range cuisinesVal.([]interface{}) {
			cuisines = append(cuisines, c.(string))
		}
	}
	
	// Parse dates
	var createdAt, updatedAt *time.Time
	if createdAtStr, ok := doc["created_at"].(string); ok && createdAtStr != "" {
		t, err := time.Parse(time.RFC3339, createdAtStr)
		if err == nil {
			createdAt = &t
		}
	}
	
	if updatedAtStr, ok := doc["updated_at"].(string); ok && updatedAtStr != "" {
		t, err := time.Parse(time.RFC3339, updatedAtStr)
		if err == nil {
			updatedAt = &t
		}
	}
	
	// Get rating values with defaults
	avgRating := 0.0
	if val, ok := doc["avg_rating"].(float64); ok {
		avgRating = val
	}
	
	ratingCount := 0
	if val, ok := doc["rating_count"].(float64); ok {
		ratingCount = int(val)
	}
	
	popularityScore := 0.0
	if val, ok := doc["popularity_score"].(float64); ok {
		popularityScore = val
	}
	
	deliveryTime := 30
	if val, ok := doc["delivery_time"].(float64); ok {
		deliveryTime = int(val)
	}
	
	// Get distance if available
	var distance *float64
	if distVal, ok := doc["distance"].(float64); ok {
		distance = &distVal
	}
	
	return RestaurantSearchResDto{
		ID:               id,
		Name:             doc["name"].(string),
		Address:          doc["address"].(string),
		Logo:             logo,
		Cover:            cover,
		ShippingFeePerKm: doc["shipping_fee_per_km"].(float64),
		AvgRating:        avgRating,
		RatingCount:      ratingCount,
		Cuisines:         cuisines,
		PopularityScore:  popularityScore,
		DeliveryTime:     deliveryTime,
		Distance:         distance,
		Status:           doc["status"].(string),
		CreatedAt:        createdAt,
		UpdatedAt:        updatedAt,
	}
}
