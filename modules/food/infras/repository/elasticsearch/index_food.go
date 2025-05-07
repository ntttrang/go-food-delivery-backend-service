package elasticsearch

import (
	"context"

	foodmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/food/model"
)

// IndexFood indexes a food document in Elasticsearch
func (r *FoodSearchRepo) IndexFood(ctx context.Context, food *foodmodel.Food) error {
	foodDoc := food.ToFoodDocument()
	return r.esClient.IndexDocument(ctx, food.Id.String(), foodDoc)
}
