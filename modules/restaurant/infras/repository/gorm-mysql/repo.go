package restaurantgormmysql

import "gorm.io/gorm"

type RestaurantRepo struct {
	db *gorm.DB
}

func NewRestaurantRepo(db *gorm.DB) *RestaurantRepo {
	return &RestaurantRepo{db: db}
}

type RestaurantFoodRepo struct {
	db *gorm.DB
}

func NewRestaurantFoodRepo(db *gorm.DB) *RestaurantFoodRepo {
	return &RestaurantFoodRepo{db: db}
}

type RestaurantLikeRepo struct {
	db *gorm.DB
}

func NewRestaurantLikeRepo(db *gorm.DB) *RestaurantLikeRepo {
	return &RestaurantLikeRepo{db: db}
}
