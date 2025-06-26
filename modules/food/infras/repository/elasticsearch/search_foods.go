package elasticsearch

import (
	"context"

	foodmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/food/model"
	"github.com/ntttrang/go-food-delivery-backend-service/modules/food/service"
	"github.com/pkg/errors"
)

// SearchFoods searches for foods based on the provided query
func (r *FoodSearchRepo) SearchFoods(ctx context.Context, req service.FoodSearchReq) (*service.FoodSearchRes, error) {
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
	res := &service.FoodSearchRes{
		Items:      items,
		Pagination: req.PagingDto,
		Facets:     processFacets(aggregations),
	}
	res.Pagination.Total = total

	return res, nil
}
