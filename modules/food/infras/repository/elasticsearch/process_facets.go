package elasticsearch

import (
	"fmt"

	foodmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/food/model"
)

// processFacets processes the aggregation results from Elasticsearch into facets
func processFacets(aggregations map[string]interface{}) foodmodel.FoodSearchFacets {
	facets := foodmodel.FoodSearchFacets{
		Categories:  []foodmodel.CategoryFacet{},
		PriceRanges: []foodmodel.PriceRangeFacet{},
		Ratings:     []foodmodel.RatingFacet{},
	}

	// If no aggregations, return empty facets
	if aggregations == nil {
		return facets
	}

	// Process category facets
	if categoryAgg, ok := aggregations["categories"].(map[string]interface{}); ok {
		if buckets, ok := categoryAgg["buckets"].([]interface{}); ok {
			for _, bucket := range buckets {
				b := bucket.(map[string]interface{})
				categoryId := b["key"].(string)
				count := int(b["doc_count"].(float64))

				// In a real implementation, you would fetch the category name from a service
				// For now, we'll just use the ID as the name
				facets.Categories = append(facets.Categories, foodmodel.CategoryFacet{
					ID:    categoryId,
					Name:  categoryId, // Ideally, this would be the category name
					Count: count,
				})
			}
		}
	}

	// Process price range facets
	if priceRangeAgg, ok := aggregations["price_ranges"].(map[string]interface{}); ok {
		if buckets, ok := priceRangeAgg["buckets"].([]interface{}); ok {
			for _, bucket := range buckets {
				b := bucket.(map[string]interface{})
				count := int(b["doc_count"].(float64))

				// Format the price range
				var rangeText string
				if from, fromOk := b["from"]; fromOk {
					if to, toOk := b["to"]; toOk {
						rangeText = fmt.Sprintf("₫%.0f - ₫%.0f", from.(float64), to.(float64))
					} else {
						rangeText = fmt.Sprintf("Over ₫%.0f", from.(float64))
					}
				} else if to, toOk := b["to"]; toOk {
					rangeText = fmt.Sprintf("Under ₫%.0f", to.(float64))
				}

				facets.PriceRanges = append(facets.PriceRanges, foodmodel.PriceRangeFacet{
					Range: rangeText,
					Count: count,
				})
			}
		}
	}

	// Process rating facets
	if ratingAgg, ok := aggregations["ratings"].(map[string]interface{}); ok {
		if buckets, ok := ratingAgg["buckets"].([]interface{}); ok {
			for _, bucket := range buckets {
				b := bucket.(map[string]interface{})
				count := int(b["doc_count"].(float64))

				// Get the rating range
				var rating float64
				if from, fromOk := b["from"]; fromOk {
					rating = from.(float64)
				}

				facets.Ratings = append(facets.Ratings, foodmodel.RatingFacet{
					Rating: rating,
					Count:  count,
				})
			}
		}
	}

	return facets
}
