package elasticsearch

import (
	"context"
	"log"

	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
)

// ReindexAllRestaurants reindexes all restaurants from the database
func (r *RestaurantSearchRepo) ReindexAllRestaurants(ctx context.Context, restaurants []restaurantmodel.RestaurantInfoDto) error {
	log.Printf("Reindexing %d restaurants", len(restaurants))
	return r.BulkIndexRestaurants(ctx, restaurants)
}
