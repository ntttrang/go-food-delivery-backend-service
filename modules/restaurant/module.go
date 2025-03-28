package restaurantmodule

import (
	"github.com/gin-gonic/gin"
	restaurantHttpgin "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/infras/controller/http-gin"
	restaurantgormmysql "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/infras/repository/gorm-mysql"
	restaurantService "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/service"
	"gorm.io/gorm"
)

func SetupRestaurantModule(db *gorm.DB, g *gin.RouterGroup) {
	restaurantRepo := restaurantgormmysql.NewRestaurantRepo(db)
	restaurantFoodRepo := restaurantgormmysql.NewRestaurantFoodRepo(db)
	restaurantLikeRepo := restaurantgormmysql.NewRestaurantLikeRepo(db)

	// TODO
	//categoryRPCClient := rpcclient.NewCategoryRPCClient(catServiceURL)

	createCmdHdl := restaurantService.NewCreateCommandHandler(restaurantRepo, restaurantFoodRepo)
	listQueryHdl := restaurantService.NewListQueryHandler(restaurantRepo, restaurantFoodRepo)
	getDetailQueryHdl := restaurantService.NewGetDetailQueryHandler(restaurantRepo)
	updateCmdHdl := restaurantService.NewUpdateRestaurantCommandHandler(restaurantRepo)
	deleteCmdHdl := restaurantService.NewDeleteCommandHandler(restaurantRepo)

	createRestaurantFavoriteCmdl := restaurantService.NewAddFavoritesCommandHandler(restaurantLikeRepo)
	deleteRestaurantFavoriteCmdl := restaurantService.NewDeleteRestaurantLikeCommandHandler(restaurantLikeRepo)

	resCtl := restaurantHttpgin.NewRestaurantHttpController(
		createCmdHdl, listQueryHdl, getDetailQueryHdl, updateCmdHdl, deleteCmdHdl,
		createRestaurantFavoriteCmdl, deleteRestaurantFavoriteCmdl)

	restaurants := g.Group("/restaurants")
	resCtl.SetupRoutes(restaurants)
}
