package elasticsearch

import (
	"context"
)

// Initialize creates the restaurant index with the proper mapping
func (r *RestaurantSearchRepo) Initialize(ctx context.Context) error {
	return r.esClient.CreateIndex(ctx, RestaurantIndexMapping)
}
