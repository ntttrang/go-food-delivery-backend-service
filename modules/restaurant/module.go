package restaurantmodule

import (
	"github.com/gin-gonic/gin"
	restaurantHttpgin "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/infras/controller/http-gin"
	restaurantgormmysql "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/infras/repository/gorm-mysql"
	restaurantService "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/service"
	shareinfras "github.com/ntttrang/go-food-delivery-backend-service/shared/infras"
)

func SetupRestaurantModule(appCtx shareinfras.IAppContext, g *gin.RouterGroup) {
	dbCtx := appCtx.DbContext()

	restaurantRepo := restaurantgormmysql.NewRestaurantRepo(dbCtx)
	restaurantFoodRepo := restaurantgormmysql.NewRestaurantFoodRepo(dbCtx)
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
	deleteRestaurantFavoriteCmdl := restaurantService.NewDeleteRestaurantLikeCommandHandler(restaurantLikeRepo)

	createCommentRestaurantCmdl := restaurantService.NewCommentRestaurantCommandHandler(restaurantRatingRepo)
	listCommentRestaurantCmdl := restaurantService.NewListRestaurantCommentsQueryHandler(restaurantRatingRepo, restaurantRepo)

	resCtl := restaurantHttpgin.NewRestaurantHttpController(
		createCmdHdl, listQueryHdl, getDetailQueryHdl, updateCmdHdl, deleteCmdHdl,
		createRestaurantFavoriteCmdl, deleteRestaurantFavoriteCmdl,
		createCommentRestaurantCmdl, listCommentRestaurantCmdl)

	restaurants := g.Group("/restaurants")
	resCtl.SetupRoutes(restaurants)
}
