package restaurantmodule

import (
	"github.com/gin-gonic/gin"
	restaurantHttpgin "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/infras/controller/http-gin"
	restaurantgormmysql "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/infras/repository/gorm-mysql"
	rpcclient "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/infras/repository/rpc-client"
	restaurantService "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/service"
	shareinfras "github.com/ntttrang/go-food-delivery-backend-service/shared/infras"
)

func SetupRestaurantModule(appCtx shareinfras.IAppContext, g *gin.RouterGroup) {
	dbCtx := appCtx.DbContext()

	foodRPCClient := rpcclient.NewFoodRPCClient(appCtx.GetConfig().FoodServiceURL)
	catRPCClient := rpcclient.NewCategoryRPCClient(appCtx.GetConfig().CatServiceURL)

	restaurantRepo := restaurantgormmysql.NewRestaurantRepo(dbCtx)
	restaurantFoodRepo := restaurantgormmysql.NewRestaurantFoodRepo(dbCtx, *foodRPCClient)
	restaurantLikeRepo := restaurantgormmysql.NewRestaurantLikeRepo(dbCtx)
	restaurantRatingRepo := restaurantgormmysql.NewRestaurantRatingRepo(dbCtx)

	// TODO
	//categoryRPCClient := rpcclient.NewCategoryRPCClient(catServiceURL)

	createCmdHdl := restaurantService.NewCreateCommandHandler(restaurantRepo, restaurantFoodRepo)
	listQueryHdl := restaurantService.NewListQueryHandler(restaurantRepo, restaurantFoodRepo)
	getDetailQueryHdl := restaurantService.NewGetDetailQueryHandler(restaurantRepo)
	updateCmdHdl := restaurantService.NewUpdateRestaurantCommandHandler(restaurantRepo)
	deleteCmdHdl := restaurantService.NewDeleteCommandHandler(restaurantRepo)

	createRestaurantFavoriteCmdl := restaurantService.NewAddFavoritesCommandHandler(restaurantLikeRepo)
	favoriteRestaurantQueryHdl := restaurantService.NewGetFavoritesRestaurantQueryHandler(restaurantRepo)

	userRPCClient := rpcclient.NewUserRPCClient(appCtx.GetConfig().UserServiceURL)

	createCommentRestaurantCmdl := restaurantService.NewCommentRestaurantCommandHandler(restaurantRatingRepo)
	listCommentRestaurantCmdl := restaurantService.NewListRestaurantCommentsQueryHandler(restaurantRatingRepo, userRPCClient)
	deleteCommentRestaurantCmdl := restaurantService.NewDeleteCommentCommandHandler(restaurantRatingRepo)

	createMenuItemCmdHdl := restaurantService.NewCreateMenuItemCommandHandler(restaurantFoodRepo)
	listMenuItemCmdHdl := restaurantService.NewListMenuItemQueryHandler(restaurantFoodRepo, foodRPCClient, catRPCClient)

	resCtl := restaurantHttpgin.NewRestaurantHttpController(
		createCmdHdl, listQueryHdl, getDetailQueryHdl, updateCmdHdl, deleteCmdHdl,
		createRestaurantFavoriteCmdl, favoriteRestaurantQueryHdl,
		createCommentRestaurantCmdl, listCommentRestaurantCmdl, deleteCommentRestaurantCmdl,
		createMenuItemCmdHdl, listMenuItemCmdHdl,
	)

	restaurants := g.Group("/restaurants")
	resCtl.SetupRoutes(restaurants)
}
