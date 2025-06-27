package elasticsearch

import (
	shareComponent "github.com/ntttrang/go-food-delivery-backend-service/shared/component"
)

type RestaurantSearchRepo struct {
	esClient *shareComponent.ElasticsearchClient
}

func NewRestaurantSearchRepo(esClient *shareComponent.ElasticsearchClient) *RestaurantSearchRepo {
	return &RestaurantSearchRepo{
		esClient: esClient,
	}
}
