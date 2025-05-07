package elasticsearch

import (
	"context"

	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	"github.com/pkg/errors"
)

// GetRestaurantById retrieves a restaurant document by ID
func (r *RestaurantSearchRepo) GetRestaurantById(ctx context.Context, id string) (*restaurantmodel.RestaurantSearchResDto, error) {
	doc, err := r.esClient.GetDocument(ctx, id)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	result := restaurantmodel.FromRestaurantDocument(doc)
	return &result, nil
}
