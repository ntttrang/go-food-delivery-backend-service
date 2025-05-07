package elasticsearch

import (
	"context"
	"log"

	foodmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/food/model"
)

// ReindexAllFoods reindexes all foods from the database
func (r *FoodSearchRepo) ReindexAllFoods(ctx context.Context, foods []foodmodel.Food) error {
	log.Printf("Reindexing %d foods", len(foods))
	return r.BulkIndexFoods(ctx, foods)
}
