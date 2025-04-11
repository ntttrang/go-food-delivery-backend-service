package restaurantgormmysql

import (
	shareinfras "github.com/ntttrang/go-food-delivery-backend-service/shared/infras"
)

// Restaurant
type RestaurantRepo struct {
	dbCtx shareinfras.IDbContext
}

func NewRestaurantRepo(dbCtx shareinfras.IDbContext) *RestaurantRepo {
	return &RestaurantRepo{dbCtx: dbCtx}
}

// Restaurant Food
type RestaurantFoodRepo struct {
	dbCtx shareinfras.IDbContext
}

func NewRestaurantFoodRepo(dbCtx shareinfras.IDbContext) *RestaurantFoodRepo {
	return &RestaurantFoodRepo{dbCtx: dbCtx}
}

// Restaurant Like
type RestaurantLikeRepo struct {
	dbCtx shareinfras.IDbContext
}

func NewRestaurantLikeRepo(dbCtx shareinfras.IDbContext) *RestaurantLikeRepo {
	return &RestaurantLikeRepo{dbCtx: dbCtx}
}

// Restaurant Rating
type RestaurantRatingRepo struct {
	dbCtx shareinfras.IDbContext
}

func NewRestaurantRatingRepo(dbCtx shareinfras.IDbContext) *RestaurantRatingRepo {
	return &RestaurantRatingRepo{dbCtx: dbCtx}
}
