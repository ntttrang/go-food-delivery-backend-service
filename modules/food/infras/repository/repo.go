package repository

import shareinfras "github.com/ntttrang/go-food-delivery-backend-service/shared/infras"

type FoodRepo struct {
	dbCtx shareinfras.IDbContext
}

func NewFoodRepo(dbCtx shareinfras.IDbContext) *FoodRepo {
	return &FoodRepo{dbCtx: dbCtx}
}

// Food Like
type FoodLikeRepo struct {
	dbCtx shareinfras.IDbContext
}

func NewFoodLikeRepo(dbCtx shareinfras.IDbContext) *FoodLikeRepo {
	return &FoodLikeRepo{dbCtx: dbCtx}
}

// Food Rating
type FoodRatingRepo struct {
	dbCtx shareinfras.IDbContext
}

func NewFoodRatingRepo(dbCtx shareinfras.IDbContext) *FoodRatingRepo {
	return &FoodRatingRepo{dbCtx: dbCtx}
}
