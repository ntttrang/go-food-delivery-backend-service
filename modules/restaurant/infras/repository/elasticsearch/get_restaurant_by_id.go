package elasticsearch

import (
	"context"

	restaurantservice "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/service"
	"github.com/pkg/errors"
)

// GetRestaurantById retrieves a restaurant document by ID
func (r *RestaurantSearchRepo) GetRestaurantById(ctx context.Context, id string) (*restaurantservice.RestaurantSearchResDto, error) {
	doc, err := r.esClient.GetDocument(ctx, id)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	modelDto := restaurantservice.FromRestaurantDocument(doc)
	result := restaurantservice.RestaurantSearchResDto{
		ID:               modelDto.ID,
		Name:             modelDto.Name,
		Address:          modelDto.Address,
		Logo:             modelDto.Logo,
		Cover:            modelDto.Cover,
		ShippingFeePerKm: modelDto.ShippingFeePerKm,
		AvgRating:        modelDto.AvgRating,
		RatingCount:      modelDto.RatingCount,
		Cuisines:         modelDto.Cuisines,
		PopularityScore:  modelDto.PopularityScore,
		DeliveryTime:     modelDto.DeliveryTime,
		Distance:         modelDto.Distance,
		Status:           modelDto.Status,
		CreatedAt:        modelDto.CreatedAt,
		UpdatedAt:        modelDto.UpdatedAt,
	}
	return &result, nil
}
