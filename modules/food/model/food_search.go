package foodmodel

import (
	"time"

	"github.com/google/uuid"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

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
	CategoryID      string   `json:"category_id"`
	RestaurantID    string   `json:"restaurant_id"`
	RestaurantName  string   `json:"restaurant_name"`
	Status          string   `json:"status"`
	Images          string   `json:"images,omitempty"`
	AvgRating       float64  `json:"avg_rating"`
	RatingCount     int      `json:"rating_count"`
	PopularityScore float64  `json:"popularity_score"`
	DeliveryTime    int      `json:"delivery_time"`
	FreeDelivery    bool     `json:"free_delivery"`
	Cuisines        []string `json:"cuisines,omitempty"` // category name
	CreatedAt       string   `json:"created_at"`
	UpdatedAt       string   `json:"updated_at"`
}

type FoodDto struct {
	Id               uuid.UUID  `json:"id"`
	RestaurantId     uuid.UUID  `json:"restaurantId"`
	RestaurantName   string     `json:"restaurantName"`
	RestaurantLat    float64    `json:"restaurantLat"`
	RestaurantLng    float64    `json:"restaurantLng"`
	ShippingFeePerKm float64    `json:"shippingFeePerKm"`
	CategoryId       uuid.UUID  `json:"categoryId"`
	CategoryName     string     `json:"categoryName"`
	Name             string     `json:"name"`
	Description      string     `json:"description"`
	Price            float64    `json:"price"`
	Images           string     `json:"images"`
	AvgPoint         float64    `json:"avgPoint"`
	CommentQty       int        `json:"commentQty"`
	Status           string     `json:"status"`
	CreatedAt        *time.Time `json:"createdAt"`
	UpdatedAt        *time.Time `json:"updatedAt"`
}

// ToFoodDocument converts a Food to a FoodDocument
func (f *FoodDto) ToFoodDocument() FoodDocument {
	var createdAt, updatedAt string
	if f.CreatedAt != nil {
		createdAt = f.CreatedAt.Format("2006-01-02T15:04:05Z")
	}
	if f.UpdatedAt != nil {
		updatedAt = f.UpdatedAt.Format("2006-01-02T15:04:05Z")
	}

	cuisines := []string{}
	if f.CategoryId != uuid.Nil {
		cuisines = append(cuisines, f.CategoryName)
	}

	freeDelivery := true
	if f.ShippingFeePerKm != 0 {
		freeDelivery = false
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
		AvgRating:       f.AvgPoint,   // This would be populated from ratings
		RatingCount:     f.CommentQty, // This would be populated from ratings
		PopularityScore: 0,            // This would be calculated based on views, orders, etc.
		DeliveryTime:    30,           // Default delivery time in minutes
		FreeDelivery:    freeDelivery, // Default to no free delivery
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
