package elasticsearch

import (
	"fmt"

	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
)

// buildRestaurantSearchQuery constructs the Elasticsearch query based on the search request
func buildRestaurantSearchQuery(req restaurantmodel.RestaurantSearchReq) map[string]interface{} {
	// Start with a bool query
	boolQuery := map[string]interface{}{
		"must":   []interface{}{},
		"filter": []interface{}{},
	}

	// Add keyword search if provided
	if req.Keyword != "" {
		boolQuery["must"] = append(boolQuery["must"].([]interface{}), map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":     req.Keyword,
				"fields":    []string{"name", "address"},
				"type":      "best_fields",
				"fuzziness": "AUTO",
			},
		})
	}

	// Add name search if provided
	if req.Name != "" {
		boolQuery["must"] = append(boolQuery["must"].([]interface{}), map[string]interface{}{
			"match": map[string]interface{}{
				"name": map[string]interface{}{
					"query": req.Name,
					"boost": 2.0,
				},
			},
		})
	}

	// Add city filter if provided
	if req.CityID != nil {
		boolQuery["filter"] = append(boolQuery["filter"].([]interface{}), map[string]interface{}{
			"term": map[string]interface{}{
				"city_id": *req.CityID,
			},
		})
	}

	// Add cuisines filter if provided
	if len(req.Cuisines) > 0 {
		boolQuery["filter"] = append(boolQuery["filter"].([]interface{}), map[string]interface{}{
			"terms": map[string]interface{}{
				"cuisines": req.Cuisines,
			},
		})
	}

	// Add rating filter if provided
	if req.Rating != nil {
		boolQuery["filter"] = append(boolQuery["filter"].([]interface{}), map[string]interface{}{
			"range": map[string]interface{}{
				"avg_rating": map[string]interface{}{
					"gte": *req.Rating,
				},
			},
		})
	}

	// Add free shipping filter if provided
	if req.FreeShipping != nil && *req.FreeShipping {
		boolQuery["filter"] = append(boolQuery["filter"].([]interface{}), map[string]interface{}{
			"term": map[string]interface{}{
				"shipping_fee_per_km": 0,
			},
		})
	}

	// Add geo distance filter if lat, lng, and radius are provided
	if req.Lat != nil && req.Lng != nil && req.Radius != nil {
		boolQuery["filter"] = append(boolQuery["filter"].([]interface{}), map[string]interface{}{
			"geo_distance": map[string]interface{}{
				"distance": fmt.Sprintf("%.1fkm", *req.Radius),
				"location": map[string]interface{}{
					"lat": *req.Lat,
					"lon": *req.Lng,
				},
			},
		})
	}

	// Add status filter (always filter for active restaurants)
	boolQuery["filter"] = append(boolQuery["filter"].([]interface{}), map[string]interface{}{
		"term": map[string]interface{}{
			"status": "ACTIVE",
		},
	})

	// If no must clauses, match all documents
	if len(boolQuery["must"].([]interface{})) == 0 {
		boolQuery["must"] = append(boolQuery["must"].([]interface{}), map[string]interface{}{
			"match_all": map[string]interface{}{},
		})
	}

	// Build sort
	sort := []interface{}{}

	if req.SortBy != "" {
		direction := "desc"
		if req.Direction == "asc" {
			direction = "asc"
		}

		// Map the sort field from the request to the Elasticsearch field
		sortField := req.SortBy
		switch req.SortBy {
		case "name":
			sortField = "name.keyword"
		case "rating":
			sortField = "avg_rating"
		case "popularity":
			sortField = "popularity_score"
		case "delivery_time":
			sortField = "delivery_time"
		case "price":
			sortField = "shipping_fee_per_km"
		case "distance":
			if req.Lat != nil && req.Lng != nil {
				// Sort by distance from the user
				sort = append(sort, map[string]interface{}{
					"_geo_distance": map[string]interface{}{
						"location": map[string]interface{}{
							"lat": *req.Lat,
							"lon": *req.Lng,
						},
						"order": direction,
						"unit":  "km",
					},
				})
				// Skip adding the regular sort field
				sortField = ""
			}
		default:
			// Default to relevance (score) sorting
			sortField = "_score"
		}

		if sortField != "" {
			sort = append(sort, map[string]interface{}{
				sortField: map[string]interface{}{
					"order": direction,
				},
			})
		}
	} else {
		// Default sort by relevance
		sort = append(sort, map[string]interface{}{
			"_score": map[string]interface{}{
				"order": "desc",
			},
		})
	}

	// Add aggregations for facets
	aggs := map[string]interface{}{
		"cuisines": map[string]interface{}{
			"terms": map[string]interface{}{
				"field": "cuisines",
				"size":  20,
			},
		},
		"ratings": map[string]interface{}{
			"range": map[string]interface{}{
				"field": "avg_rating",
				"ranges": []interface{}{
					map[string]interface{}{"from": 4, "to": 5},
					map[string]interface{}{"from": 3, "to": 4},
					map[string]interface{}{"from": 2, "to": 3},
					map[string]interface{}{"from": 1, "to": 2},
					map[string]interface{}{"from": 0, "to": 1},
				},
			},
		},
		"delivery_time": map[string]interface{}{
			"range": map[string]interface{}{
				"field": "delivery_time",
				"ranges": []interface{}{
					map[string]interface{}{"to": 15},
					map[string]interface{}{"from": 15, "to": 30},
					map[string]interface{}{"from": 30, "to": 45},
					map[string]interface{}{"from": 45, "to": 60},
					map[string]interface{}{"from": 60},
				},
			},
		},
	}

	// Construct the final query
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": boolQuery,
		},
		"sort":             sort,
		"aggs":             aggs,
		"track_total_hits": true,
	}

	return query
}
