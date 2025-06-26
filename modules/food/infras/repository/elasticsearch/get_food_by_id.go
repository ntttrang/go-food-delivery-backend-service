package elasticsearch

import (
	"context"

	foodmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/food/model"
	"github.com/pkg/errors"
)

// GetFoodById retrieves a food document by ID
func (r *FoodSearchRepo) GetFoodById(ctx context.Context, id string) (*foodmodel.FoodSearchResDto, error) {
	doc, err := r.esClient.GetDocument(ctx, id)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	result := foodmodel.FromFoodDocument(doc)
	return &result, nil
}
