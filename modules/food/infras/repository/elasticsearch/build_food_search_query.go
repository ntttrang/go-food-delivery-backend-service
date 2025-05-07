package elasticsearch

import (
	foodmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/food/model"
)

// buildFoodSearchQuery constructs the Elasticsearch query based on the search request
func buildFoodSearchQuery(req foodmodel.FoodSearchReq) map[string]interface{} {
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
				"fields":    []string{"name", "cuisine", "city"},
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

	// Add description search if provided
	if req.Description != "" {
		boolQuery["must"] = append(boolQuery["must"].([]interface{}), map[string]interface{}{
			"match": map[string]interface{}{
				"description": req.Description,
			},
		})
	}

	// Add category filter if provided
	if len(req.CategoryIds) > 0 {
		boolQuery["filter"] = append(boolQuery["filter"].([]interface{}), map[string]interface{}{
			"terms": map[string]interface{}{
				"category_id": req.CategoryIds,
			},
		})
	}

	// Add price range filter if provided
	if req.PriceMin != nil || req.PriceMax != nil {
		rangeQuery := map[string]interface{}{}

		if req.PriceMin != nil {
			rangeQuery["gte"] = *req.PriceMin
		}

		if req.PriceMax != nil {
			rangeQuery["lte"] = *req.PriceMax
		}

		boolQuery["filter"] = append(boolQuery["filter"].([]interface{}), map[string]interface{}{
			"range": map[string]interface{}{
				"price": rangeQuery,
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

	// Add status filter (always filter for active foods)
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
		case "price":
			sortField = "price"
		case "rating":
			sortField = "avg_rating"
		case "created_at":
			sortField = "created_at"
		case "popularity":
			// Sort by popularity score (could be based on views, orders, etc.)
			sortField = "popularity_score"
		case "delivery_time":
			// Sort by delivery time
			sortField = "delivery_time"
		default:
			// Default to relevance (score) sorting
			sortField = "_score"
		}

		sort = append(sort, map[string]interface{}{
			sortField: map[string]interface{}{
				"order": direction,
			},
		})
	} else {
		// Default sort by relevance (score) and then by created_at
		sort = append(sort, map[string]interface{}{"_score": map[string]interface{}{"order": "desc"}})
		sort = append(sort, map[string]interface{}{"created_at": map[string]interface{}{"order": "desc"}})
	}

	// Add aggregations for facets
	aggs := map[string]interface{}{
		"categories": map[string]interface{}{
			"terms": map[string]interface{}{
				"field": "category_id",
				"size":  10,
			},
		},
		"price_ranges": map[string]interface{}{
			"range": map[string]interface{}{
				"field": "price",
				"ranges": []interface{}{
					map[string]interface{}{"to": 50000},
					map[string]interface{}{"from": 50000, "to": 100000},
					map[string]interface{}{"from": 100000, "to": 200000},
					map[string]interface{}{"from": 200000, "to": 500000},
					map[string]interface{}{"from": 500000},
				},
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
