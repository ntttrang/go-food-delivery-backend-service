package restaurantmodel

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	sharedmodel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
	"gorm.io/gorm"
)

type Restaurant struct {
	Id               uuid.UUID       `gorm:"column:id;"`
	OwnerId          uuid.UUID       `gorm:"column:owner_id;"`
	Name             string          `gorm:"column:name;"`
	Addr             string          `gorm:"column:addr;"`
	CityId           int             `gorm:"column:city_id;"`
	Lat              float64         `gorm:"column:lat;"`
	Lng              float64         `gorm:"column:lng;"`
	Cover            json.RawMessage `gorm:"column:cover;"` // json
	Logo             json.RawMessage `gorm:"column:logo;"`  // json
	ShippingFeePerKm float64         `gorm:"column:shipping_fee_per_km;"`
	Status           string          `gorm:"column:status;"`
	sharedmodel.DateDto
}

func (Restaurant) TableName() string {
	return "restaurants"
}

func (r *Restaurant) BeforeCreate(tx *gorm.DB) {
	//currentTime := time.Now().UTC()
	currentTime := time.Date(2021, time.March, 28, 12, 31, 0, 0, time.UTC)
	r.CreatedAt = &currentTime
	r.UpdatedAt = &currentTime
}

// ToRestaurantDocument converts a Restaurant to an Elasticsearch document
func (r *Restaurant) ToRestaurantDocument() map[string]any {
	// Create a base document with all the fields
	doc := map[string]any{
		"id":                  r.Id.String(),
		"name":                r.Name,
		"address":             r.Addr,
		"city_id":             r.CityId,
		"lat":                 r.Lat,
		"lng":                 r.Lng,
		"shipping_fee_per_km": r.ShippingFeePerKm,
		"status":              r.Status,
	}

	// Add timestamps if available
	if r.CreatedAt != nil {
		doc["created_at"] = r.CreatedAt.Format(time.RFC3339)
	}
	if r.UpdatedAt != nil {
		doc["updated_at"] = r.UpdatedAt.Format(time.RFC3339)
	}

	// Add complex fields (logo, cover) if available
	if len(r.Logo) > 0 {
		var logoMap map[string]any
		if err := json.Unmarshal(r.Logo, &logoMap); err == nil {
			doc["logo"] = logoMap
		}
	}
	if len(r.Cover) > 0 {
		var coverMap map[string]any
		if err := json.Unmarshal(r.Cover, &coverMap); err == nil {
			doc["cover"] = coverMap
		}
	}

	// Default values for search-specific fields
	doc["avg_rating"] = 0.0
	doc["rating_count"] = 0
	doc["popularity_score"] = 0.0
	doc["delivery_time"] = 30 // Default delivery time in minutes
	doc["cuisines"] = []string{}

	return doc
}
