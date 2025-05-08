package restaurantmodule

import (
	"log"

	"github.com/gin-gonic/gin"
	restaurantHttpgin "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/infras/controller/http-gin"
	elasticsearchrepo "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/infras/repository/elasticsearch"
	restaurantgormmysql "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/infras/repository/gorm-mysql"
	rpcclient "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/infras/repository/rpc-client"
	restaurantService "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/service"
	shareComponent "github.com/ntttrang/go-food-delivery-backend-service/shared/component"
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

	userRPCClient := rpcclient.NewUserRPCClient(appCtx.GetConfig().UserServiceURL)

	esClient, err := shareComponent.NewElasticsearchClient(appCtx.GetConfig().ElasticSearch)
	if err != nil {
		log.Printf("Elasticsearch initialization failed: %v. Search functionality will be disabled.", err)
	}
	esRestaurantRepo := elasticsearchrepo.NewRestaurantSearchRepo(esClient)

	// Setup basic handlers
	createCmdHdl := restaurantService.NewCreateCommandHandler(restaurantRepo, restaurantFoodRepo)
	listQueryHdl := restaurantService.NewListQueryHandler(restaurantRepo, restaurantFoodRepo)
	getDetailQueryHdl := restaurantService.NewGetDetailQueryHandler(restaurantRepo)
	updateCmdHdl := restaurantService.NewUpdateRestaurantCommandHandler(restaurantRepo)
	deleteCmdHdl := restaurantService.NewDeleteCommandHandler(restaurantRepo)

	// Favorite restaurant
	createRestaurantFavoriteCmdl := restaurantService.NewAddFavoritesCommandHandler(restaurantLikeRepo)
	favoriteRestaurantQueryHdl := restaurantService.NewGetFavoritesRestaurantQueryHandler(restaurantRepo)

	// restaurant comment
	createCommentRestaurantCmdl := restaurantService.NewCommentRestaurantCommandHandler(restaurantRatingRepo)
	listCommentRestaurantCmdl := restaurantService.NewListRestaurantCommentsQueryHandler(restaurantRatingRepo, userRPCClient)
	deleteCommentRestaurantCmdl := restaurantService.NewDeleteCommentCommandHandler(restaurantRatingRepo)

	createMenuItemCmdHdl := restaurantService.NewCreateMenuItemCommandHandler(restaurantFoodRepo)
	listMenuItemCmdHdl := restaurantService.NewListMenuItemQueryHandler(restaurantFoodRepo, foodRPCClient, catRPCClient)
	deleteMenuItemCmdHdl := restaurantService.NewDeleteMenuItemCommandHandler(restaurantFoodRepo)

	// Setup Elasticsearch
	searchRestaurantQueryHandler := restaurantService.NewSearchRestaurantQueryHandler(esRestaurantRepo)
	syncRestaurantByIdCmdHdl := restaurantService.NewSyncRestaurantByIdCommandHandler(restaurantRepo, esRestaurantRepo)
	syncRestaurantIndexCmdHdl := restaurantService.NewSyncRestaurantIndexCommandHandler(restaurantRepo, esRestaurantRepo)

	resCtl := restaurantHttpgin.NewRestaurantHttpController(
		createCmdHdl, listQueryHdl, getDetailQueryHdl, updateCmdHdl, deleteCmdHdl,
		createRestaurantFavoriteCmdl, favoriteRestaurantQueryHdl,
		createCommentRestaurantCmdl, listCommentRestaurantCmdl, deleteCommentRestaurantCmdl,
		createMenuItemCmdHdl, listMenuItemCmdHdl, deleteMenuItemCmdHdl,
		searchRestaurantQueryHandler, syncRestaurantByIdCmdHdl, syncRestaurantIndexCmdHdl,
	)

	restaurants := g.Group("/restaurants")
	resCtl.SetupRoutes(restaurants)
}
