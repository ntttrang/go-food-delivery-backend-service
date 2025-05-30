package elasticsearch

import (
	"context"

	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
)

// BulkIndexRestaurants indexes multiple restaurants in a single request
func (r *RestaurantSearchRepo) BulkIndexRestaurants(ctx context.Context, restaurants []restaurantmodel.RestaurantInfoDto) error {
	if len(restaurants) == 0 {
		return nil
	}

	// Prepare documents for bulk indexing
	documents := make(map[string]any, len(restaurants))
	for _, restaurant := range restaurants {
		restaurantDoc := restaurant.ToRestaurantDocument()
		documents[restaurant.Id.String()] = restaurantDoc
	}

	return r.esClient.BulkIndex(ctx, documents)
}
