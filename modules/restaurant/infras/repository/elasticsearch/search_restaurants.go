package elasticsearch

import (
	"context"

	restaurantservice "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/service"
	"github.com/pkg/errors"
)

// SearchRestaurants searches for restaurants based on the provided query
func (r *RestaurantSearchRepo) SearchRestaurants(ctx context.Context, req restaurantservice.RestaurantSearchReq) (*restaurantservice.RestaurantSearchRes, error) {
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
	items := make([]restaurantservice.RestaurantSearchResDto, len(results))
	for i, result := range results {
		// Convert from model to service DTO
		modelDto := restaurantservice.FromRestaurantDocument(result)
		items[i] = restaurantservice.RestaurantSearchResDto{
			ID:               modelDto.ID,
			Name:             modelDto.Name,
			Address:          modelDto.Address,
			Logo:             modelDto.Logo,
			Cover:            modelDto.Cover,
			ShippingFeePerKm: modelDto.ShippingFeePerKm,
			AvgRating:        modelDto.AvgRating,
			RatingCount:      modelDto.RatingCount,
			Cuisines:         modelDto.Cuisines,
			PopularityScore:  modelDto.PopularityScore,
			DeliveryTime:     modelDto.DeliveryTime,
			Distance:         modelDto.Distance,
			Status:           modelDto.Status,
			CreatedAt:        modelDto.CreatedAt,
			UpdatedAt:        modelDto.UpdatedAt,
		}
	}

	// Create response
	res := &restaurantservice.RestaurantSearchRes{
		Items:      items,
		Pagination: req.PagingDto,
		Facets:     processRestaurantFacets(aggregations),
	}
	res.Pagination.Total = total

	return res, nil
}
