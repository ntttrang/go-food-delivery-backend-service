package elasticsearch

import (
	"context"

	"github.com/google/uuid"
)

// DeleteFood deletes a food document from Elasticsearch
func (r *FoodSearchRepo) DeleteFood(ctx context.Context, id uuid.UUID) error {
	return r.esClient.DeleteDocument(ctx, id.String())
}
