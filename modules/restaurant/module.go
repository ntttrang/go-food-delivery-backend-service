package restaurantmodule

import (
	"github.com/gin-gonic/gin"
	restaurantHttpgin "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/infras/controller/http-gin"
	restaurantgormmysql "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/infras/respository/gorm-mysql"
	restaurantService "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/internal/service"
	"gorm.io/gorm"
)

func SetupRestaurantModule(db *gorm.DB, g *gin.RouterGroup) {
	restaurantRepo := restaurantgormmysql.NewRestaurantRepo(db)
	restaurantFoodRepo := restaurantgormmysql.NewRestaurantFoodRepo(db)

	resService := restaurantService.NewRestaurantService(restaurantRepo, restaurantFoodRepo)

	resCtl := restaurantHttpgin.NewRestaurantHttpController(resService)

	resCtl.SetupRoutes(g)
}
