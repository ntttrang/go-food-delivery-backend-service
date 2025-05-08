package elasticsearch

import (
	"context"

	foodmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/food/model"
)

// BulkIndexFoods indexes multiple foods in a single request
func (r *FoodSearchRepo) BulkIndexFoods(ctx context.Context, foods []foodmodel.Food) error {
	if len(foods) == 0 {
		return nil
	}

	// Prepare documents for bulk indexing
	documents := make(map[string]interface{}, len(foods))
	for _, food := range foods {
		foodDoc := food.ToFoodDocument()
		documents[food.Id.String()] = foodDoc
	}

	return r.esClient.BulkIndex(ctx, documents)
}
