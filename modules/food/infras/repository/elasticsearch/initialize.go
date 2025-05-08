package elasticsearch

import (
	"context"
)

// Initialize creates the food index with the proper mapping
func (r *FoodSearchRepo) Initialize(ctx context.Context) error {
	return r.esClient.CreateIndex(ctx, FoodIndexMapping)
}
