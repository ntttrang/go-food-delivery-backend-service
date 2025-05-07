package elasticsearch

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	shareComponent "github.com/ntttrang/go-food-delivery-backend-service/shared/component"
	"github.com/pkg/errors"
)

type RestaurantSearchRepo struct {
	esClient *shareComponent.ElasticsearchClient
}

func NewRestaurantSearchRepo(esClient *shareComponent.ElasticsearchClient) *RestaurantSearchRepo {
	return &RestaurantSearchRepo{
		esClient: esClient,
	}
}

// Initialize creates the restaurant index with the proper mapping
func (r *RestaurantSearchRepo) Initialize(ctx context.Context) error {
	return r.esClient.CreateIndex(ctx, RestaurantIndexMapping)
}

// IndexRestaurant indexes a restaurant document in Elasticsearch
func (r *RestaurantSearchRepo) IndexRestaurant(ctx context.Context, restaurant *restaurantmodel.Restaurant) error {
	restaurantDoc := restaurant.ToRestaurantDocument()
	return r.esClient.IndexDocument(ctx, restaurant.Id.String(), restaurantDoc)
}

// DeleteRestaurant deletes a restaurant document from Elasticsearch
func (r *RestaurantSearchRepo) DeleteRestaurant(ctx context.Context, id uuid.UUID) error {
	return r.esClient.DeleteDocument(ctx, id.String())
}

// SearchRestaurants searches for restaurants based on the provided query
func (r *RestaurantSearchRepo) SearchRestaurants(ctx context.Context, req restaurantmodel.RestaurantSearchReq) (*restaurantmodel.RestaurantSearchRes, error) {
	// Build the Elasticsearch query
	query := buildRestaurantSearchQuery(req)

	// Calculate from based on pagination
	from := (req.Page - 1) * req.Limit

	// Execute the search
	results, total, aggregations, err := r.esClient.Search(ctx, query, from, req.Limit)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// Convert results to RestaurantSearchResDto
	items := make([]restaurantmodel.RestaurantSearchResDto, len(results))
	for i, result := range results {
		items[i] = restaurantmodel.FromRestaurantDocument(result)
	}

	// Create response
	res := &restaurantmodel.RestaurantSearchRes{
		Items:      items,
		Pagination: req.PagingDto,
		Facets:     processRestaurantFacets(aggregations),
	}
	res.Pagination.Total = total

	return res, nil
}

// BulkIndexRestaurants indexes multiple restaurants in a single request
func (r *RestaurantSearchRepo) BulkIndexRestaurants(ctx context.Context, restaurants []restaurantmodel.Restaurant) error {
	if len(restaurants) == 0 {
		return nil
	}

	// Prepare documents for bulk indexing
	documents := make(map[string]interface{}, len(restaurants))
	for _, restaurant := range restaurants {
		restaurantDoc := restaurant.ToRestaurantDocument()
		documents[restaurant.Id.String()] = restaurantDoc
	}

	return r.esClient.BulkIndex(ctx, documents)
}

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

// GetRestaurantById retrieves a restaurant document by ID
func (r *RestaurantSearchRepo) GetRestaurantById(ctx context.Context, id string) (*restaurantmodel.RestaurantSearchResDto, error) {
	doc, err := r.esClient.GetDocument(ctx, id)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	result := restaurantmodel.FromRestaurantDocument(doc)
	return &result, nil
}

// ReindexAllRestaurants reindexes all restaurants from the database
func (r *RestaurantSearchRepo) ReindexAllRestaurants(ctx context.Context, restaurants []restaurantmodel.Restaurant) error {
	log.Printf("Reindexing %d restaurants", len(restaurants))
	return r.BulkIndexRestaurants(ctx, restaurants)
}

// processRestaurantFacets processes the aggregation results from Elasticsearch into facets
func processRestaurantFacets(aggregations map[string]interface{}) restaurantmodel.RestaurantSearchFacets {
	facets := restaurantmodel.RestaurantSearchFacets{
		Cuisines:     []restaurantmodel.CuisineFacet{},
		PriceRanges:  []restaurantmodel.PriceRangeFacet{},
		Ratings:      []restaurantmodel.RatingFacet{},
		DeliveryTime: []restaurantmodel.TimeFacet{},
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

				facets.Cuisines = append(facets.Cuisines, restaurantmodel.CuisineFacet{
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

				facets.Ratings = append(facets.Ratings, restaurantmodel.RatingFacet{
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

				facets.DeliveryTime = append(facets.DeliveryTime, restaurantmodel.TimeFacet{
					From:  from,
					To:    to,
					Count: count,
				})
			}
		}
	}

	return facets
}
