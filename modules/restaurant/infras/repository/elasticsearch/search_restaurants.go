package elasticsearch

import (
	"context"

	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	"github.com/pkg/errors"
)

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
