package elasticsearch

import (
	shareComponent "github.com/ntttrang/go-food-delivery-backend-service/shared/component"
)

type FoodSearchRepo struct {
	esClient *shareComponent.ElasticsearchClient
}

func NewFoodSearchRepo(esClient *shareComponent.ElasticsearchClient) *FoodSearchRepo {
	return &FoodSearchRepo{
		esClient: esClient,
	}
}
