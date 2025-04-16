package repository

import shareinfras "github.com/ntttrang/go-food-delivery-backend-service/shared/infras"

type FoodRepo struct {
	dbCtx shareinfras.IDbContext
}

func NewFoodRepo(dbCtx shareinfras.IDbContext) *FoodRepo {
	return &FoodRepo{dbCtx: dbCtx}
}
