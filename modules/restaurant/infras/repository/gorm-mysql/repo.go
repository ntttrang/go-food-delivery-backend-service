package restaurantgormmysql

import "gorm.io/gorm"

// Restaurant
type RestaurantRepo struct {
	db *gorm.DB
}

func NewRestaurantRepo(db *gorm.DB) *RestaurantRepo {
	return &RestaurantRepo{db: db}
}

// Restaurant Food
type RestaurantFoodRepo struct {
	db *gorm.DB
}

func NewRestaurantFoodRepo(db *gorm.DB) *RestaurantFoodRepo {
	return &RestaurantFoodRepo{db: db}
}

// Restaurant Like
type RestaurantLikeRepo struct {
	db *gorm.DB
}

func NewRestaurantLikeRepo(db *gorm.DB) *RestaurantLikeRepo {
	return &RestaurantLikeRepo{db: db}
}

// Restaurant Rating
type RestaurantRatingRepo struct {
	db *gorm.DB
}

func NewRestaurantRatingRepo(db *gorm.DB) *RestaurantRatingRepo {
	return &RestaurantRatingRepo{db: db}
}
