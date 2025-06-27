package elasticsearch

import (
	"context"

	"github.com/google/uuid"
)

// DeleteRestaurant deletes a restaurant document from Elasticsearch
func (r *RestaurantSearchRepo) DeleteRestaurant(ctx context.Context, id uuid.UUID) error {
	return r.esClient.DeleteDocument(ctx, id.String())
}
