package elasticsearch

import (
	"fmt"

	restaurantservice "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/service"
)

// buildRestaurantSearchQuery constructs the Elasticsearch query based on the search request
func buildRestaurantSearchQuery(req restaurantservice.RestaurantSearchReq) map[string]any {
	// Start with a bool query
	boolQuery := map[string]any{
		"must":   []any{},
		"filter": []any{},
	}

	// Add keyword search if provided
	if req.Keyword != "" {
		boolQuery["must"] = append(boolQuery["must"].([]any), map[string]any{
			"multi_match": map[string]any{
				"query":     req.Keyword,
				"fields":    []string{"name", "address"},
				"type":      "best_fields",
				"fuzziness": "AUTO",
			},
		})
	}

	// Add name search if provided
	if req.Name != "" {
		boolQuery["must"] = append(boolQuery["must"].([]any), map[string]any{
			"match": map[string]any{
				"name": map[string]any{
					"query": req.Name,
					"boost": 2.0,
				},
			},
		})
	}

	// Add city filter if provided
	if req.CityID != nil {
		boolQuery["filter"] = append(boolQuery["filter"].([]any), map[string]any{
			"term": map[string]any{
				"city_id": *req.CityID,
			},
		})
	}

	// Add cuisines filter if provided
	if len(req.Cuisines) > 0 {
		boolQuery["filter"] = append(boolQuery["filter"].([]any), map[string]any{
			"terms": map[string]any{
				"cuisines": req.Cuisines,
			},
		})
	}

	// Add rating filter if provided
	if req.Rating != nil {
		boolQuery["filter"] = append(boolQuery["filter"].([]any), map[string]any{
			"range": map[string]any{
				"avg_rating": map[string]any{
					"gte": *req.Rating,
				},
			},
		})
	}

	// Add free shipping filter if provided
	if req.FreeShipping != nil && *req.FreeShipping {
		boolQuery["filter"] = append(boolQuery["filter"].([]any), map[string]any{
			"term": map[string]any{
				"shipping_fee_per_km": 0,
			},
		})
	}

	// Add geo distance filter if lat, lng, and radius are provided
	if req.Lat != nil && req.Lng != nil && req.Radius != nil {
		boolQuery["filter"] = append(boolQuery["filter"].([]any), map[string]any{
			"geo_distance": map[string]any{
				"distance": fmt.Sprintf("%.1fkm", *req.Radius),
				"location": map[string]any{
					"lat": *req.Lat,
					"lon": *req.Lng,
				},
			},
		})
	}

	// Add status filter (always filter for active restaurants)
	boolQuery["filter"] = append(boolQuery["filter"].([]any), map[string]any{
		"term": map[string]any{
			"status": "ACTIVE",
		},
	})

	// If no must clauses, match all documents
	if len(boolQuery["must"].([]any)) == 0 {
		boolQuery["must"] = append(boolQuery["must"].([]any), map[string]any{
			"match_all": map[string]any{},
		})
	}

	// Build sort
	sort := []any{}

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
				sort = append(sort, map[string]any{
					"_geo_distance": map[string]any{
						"location": map[string]any{
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
			sort = append(sort, map[string]any{
				sortField: map[string]any{
					"order": direction,
				},
			})
		}
	} else {
		// Default sort by relevance
		sort = append(sort, map[string]any{
			"_score": map[string]any{
				"order": "desc",
			},
		})
	}

	// Add aggregations for facets
	aggs := map[string]any{
		"cuisines": map[string]any{
			"terms": map[string]any{
				"field": "cuisines",
				"size":  20,
			},
		},
		"ratings": map[string]any{
			"range": map[string]any{
				"field": "avg_rating",
				"ranges": []any{
					map[string]any{"from": 4, "to": 5},
					map[string]any{"from": 3, "to": 4},
					map[string]any{"from": 2, "to": 3},
					map[string]any{"from": 1, "to": 2},
					map[string]any{"from": 0, "to": 1},
				},
			},
		},
		"delivery_time": map[string]any{
			"range": map[string]any{
				"field": "delivery_time",
				"ranges": []any{
					map[string]any{"to": 15},
					map[string]any{"from": 15, "to": 30},
					map[string]any{"from": 30, "to": 45},
					map[string]any{"from": 45, "to": 60},
					map[string]any{"from": 60},
				},
			},
		},
	}

	// Construct the final query
	query := map[string]any{
		"query": map[string]any{
			"bool": boolQuery,
		},
		"sort":             sort,
		"aggs":             aggs,
		"track_total_hits": true,
	}

	return query
}
