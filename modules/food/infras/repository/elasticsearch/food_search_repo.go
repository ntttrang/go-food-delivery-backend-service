package elasticsearch

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	foodmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/food/model"
	shareComponent "github.com/ntttrang/go-food-delivery-backend-service/shared/component"
	"github.com/pkg/errors"
)

type FoodSearchRepo struct {
	esClient *shareComponent.ElasticsearchClient
}

func NewFoodSearchRepo(esClient *shareComponent.ElasticsearchClient) *FoodSearchRepo {
	return &FoodSearchRepo{
		esClient: esClient,
	}
}

// Initialize creates the food index with the proper mapping
func (r *FoodSearchRepo) Initialize(ctx context.Context) error {
	return r.esClient.CreateIndex(ctx, FoodIndexMapping)
}

// IndexFood indexes a food document in Elasticsearch
func (r *FoodSearchRepo) IndexFood(ctx context.Context, food *foodmodel.Food) error {
	foodDoc := food.ToFoodDocument()
	return r.esClient.IndexDocument(ctx, food.Id.String(), foodDoc)
}

// DeleteFood deletes a food document from Elasticsearch
func (r *FoodSearchRepo) DeleteFood(ctx context.Context, id uuid.UUID) error {
	return r.esClient.DeleteDocument(ctx, id.String())
}

// SearchFoods searches for foods based on the provided query
func (r *FoodSearchRepo) SearchFoods(ctx context.Context, req foodmodel.FoodSearchReq) (*foodmodel.FoodSearchRes, error) {
	// Build the Elasticsearch query
	query := buildFoodSearchQuery(req)

	// Calculate from based on pagination
	from := (req.Page - 1) * req.Limit

	// Execute the search
	results, total, aggregations, err := r.esClient.Search(ctx, query, from, req.Limit)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// Convert results to FoodSearchResDto
	items := make([]foodmodel.FoodSearchResDto, len(results))
	for i, result := range results {
		items[i] = foodmodel.FromFoodDocument(result)
	}

	// Create response
	res := &foodmodel.FoodSearchRes{
		Items:      items,
		Pagination: req.PagingDto,
		Facets:     processFacets(aggregations),
	}
	res.Pagination.Total = total

	return res, nil
}

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

// BulkIndexFoods indexes multiple foods in a single request
func (r *FoodSearchRepo) BulkIndexFoods(ctx context.Context, foods []foodmodel.Food) error {
	if len(foods) == 0 {
		return nil
	}

	// Prepare documents for bulk indexing
	documents := make(map[string]interface{}, len(foods))
	for _, food := range foods {
		foodDoc := food.ToFoodDocument()
		documents[food.Id.String()] = foodDoc
	}

	return r.esClient.BulkIndex(ctx, documents)
}

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

// GetFoodById retrieves a food document by ID
func (r *FoodSearchRepo) GetFoodById(ctx context.Context, id string) (*foodmodel.FoodSearchResDto, error) {
	doc, err := r.esClient.GetDocument(ctx, id)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	result := foodmodel.FromFoodDocument(doc)
	return &result, nil
}

// ReindexAllFoods reindexes all foods from the database
func (r *FoodSearchRepo) ReindexAllFoods(ctx context.Context, foods []foodmodel.Food) error {
	log.Printf("Reindexing %d foods", len(foods))
	return r.BulkIndexFoods(ctx, foods)
}
