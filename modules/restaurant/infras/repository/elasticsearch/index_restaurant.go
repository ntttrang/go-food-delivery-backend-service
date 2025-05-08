package elasticsearch

import (
	"context"

	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
)

// IndexRestaurant indexes a restaurant document in Elasticsearch
func (r *RestaurantSearchRepo) IndexRestaurant(ctx context.Context, restaurant *restaurantmodel.Restaurant) error {
	restaurantDoc := restaurant.ToRestaurantDocument()
	return r.esClient.IndexDocument(ctx, restaurant.Id.String(), restaurantDoc)
}
