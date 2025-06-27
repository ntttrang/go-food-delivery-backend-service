package elasticsearch

import (
	restaurantservice "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/service"
)

// processRestaurantFacets processes the aggregation results from Elasticsearch into facets
func processRestaurantFacets(aggregations map[string]interface{}) restaurantservice.RestaurantSearchFacets {
	facets := restaurantservice.RestaurantSearchFacets{
		Cuisines:     []restaurantservice.CuisineFacet{},
		PriceRanges:  []restaurantservice.PriceRangeFacet{},
		Ratings:      []restaurantservice.RatingFacet{},
		DeliveryTime: []restaurantservice.TimeFacet{},
	}

	// If no aggregations, return empty facets
	if aggregations == nil {
		return facets
	}

	// Process cuisine facets
	if cuisineAgg, ok := aggregations["cuisines"].(map[string]interface{}); ok {
		if buckets, ok := cuisineAgg["buckets"].([]interface{}); ok {
			for _, bucket := range buckets {
				b := bucket.(map[string]interface{})
				value := b["key"].(string)
				count := int(b["doc_count"].(float64))

				facets.Cuisines = append(facets.Cuisines, restaurantservice.CuisineFacet{
					Value: value,
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

				from := 0.0
				to := 5.0

				if fromVal, ok := b["from"]; ok {
					from = fromVal.(float64)
				}

				if toVal, ok := b["to"]; ok {
					to = toVal.(float64)
				}

				facets.Ratings = append(facets.Ratings, restaurantservice.RatingFacet{
					From:  from,
					To:    to,
					Count: count,
				})
			}
		}
	}

	// Process delivery time facets
	if timeAgg, ok := aggregations["delivery_time"].(map[string]interface{}); ok {
		if buckets, ok := timeAgg["buckets"].([]interface{}); ok {
			for _, bucket := range buckets {
				b := bucket.(map[string]interface{})
				count := int(b["doc_count"].(float64))

				from := 0
				to := 0

				if fromVal, ok := b["from"]; ok {
					from = int(fromVal.(float64))
				}

				if toVal, ok := b["to"]; ok {
					to = int(toVal.(float64))
				}

				facets.DeliveryTime = append(facets.DeliveryTime, restaurantservice.TimeFacet{
					From:  from,
					To:    to,
					Count: count,
				})
			}
		}
	}

	return facets
}
