package foodmodule

import (
	"github.com/gin-gonic/gin"
	foodHttpgin "github.com/ntttrang/go-food-delivery-backend-service/modules/food/infras/controller/http-gin"
	elasticsearch "github.com/ntttrang/go-food-delivery-backend-service/modules/food/infras/repository/elasticsearch"
	repo "github.com/ntttrang/go-food-delivery-backend-service/modules/food/infras/repository/gorm-mysql"
	rpcclient "github.com/ntttrang/go-food-delivery-backend-service/modules/food/infras/repository/rpc-client"
	foodService "github.com/ntttrang/go-food-delivery-backend-service/modules/food/service"
	shareComponent "github.com/ntttrang/go-food-delivery-backend-service/shared/component"
	shareinfras "github.com/ntttrang/go-food-delivery-backend-service/shared/infras"
)

func SetupFoodModule(appCtx shareinfras.IAppContext, g *gin.RouterGroup) {
	dbCtx := appCtx.DbContext()

	// Setup controller
	foodRepo := repo.NewFoodRepo(dbCtx)
	foodLikeRepo := repo.NewFoodLikeRepo(dbCtx)
	foodRatingRepo := repo.NewFoodRatingRepo(dbCtx)

	createCmdHdl := foodService.NewCreateCommandHandler(foodRepo)
	listCmdHdl := foodService.NewListCommandHandler(foodRepo)
	getDetailCmdHdl := foodService.NewGetDetailQueryHandler(foodRepo)
	updateCmdHdl := foodService.NewUpdateCommandHandler(foodRepo)
	deleteCmdHdl := foodService.NewDeleteByIdCommandHandler(foodRepo)

	createFoodFavoriteCmdl := foodService.NewAddFavoritesCommandHandler(foodLikeRepo)
	favoriteFoodQueryHdl := foodService.NewGetFavoritesFoodQueryHandler(foodRepo)

	userRPCClient := rpcclient.NewUserRPCClient(appCtx.GetConfig().UserServiceURL)

	createCommentFoodCmdl := foodService.NewCommentFoodCommandHandler(foodRatingRepo)
	listCommentFoodCmdl := foodService.NewListFoodCommentsQueryHandler(foodRatingRepo, userRPCClient)
	deleteCommentFoodCmdl := foodService.NewDeleteCommentCommandHandler(foodRatingRepo)

	// Setup Elasticsearch
	// Try to initialize Elasticsearch client
	esClient, err := shareComponent.NewElasticsearchClient(appCtx.GetConfig().ElasticSearch)
	if err != nil {
		panic(err)
	}
	foodSearchRepo := elasticsearch.NewFoodSearchRepo(esClient)
	rpcCategoryRepo := rpcclient.NewCategoryRPCClient(appCtx.GetConfig().FoodServiceURL)
	rpcRestaurantRepo := rpcclient.NewRestaurantRPCClient(appCtx.GetConfig().RestaurantServiceURL)

	searchFoodQueryHdl := foodService.NewSearchFoodQueryHandler(foodSearchRepo)
	syncFoodByIdCmdHdl := foodService.NewSyncFoodByIdCommandHandler(foodRepo, foodSearchRepo, rpcRestaurantRepo, rpcCategoryRepo)
	syncFoodIndexCmdHdl := foodService.NewSyncFoodIndexCommandHandler(foodRepo, foodSearchRepo, rpcRestaurantRepo, rpcCategoryRepo)

	foodCtrl := foodHttpgin.NewFoodHttpController(
		createCmdHdl, listCmdHdl, getDetailCmdHdl, updateCmdHdl, deleteCmdHdl, foodRepo,
		createFoodFavoriteCmdl, favoriteFoodQueryHdl,
		createCommentFoodCmdl, listCommentFoodCmdl, deleteCommentFoodCmdl,
		searchFoodQueryHdl, syncFoodByIdCmdHdl, syncFoodIndexCmdHdl,
	)

	// Setup router
	// RPC
	g.POST("/rpc/foods/find-by-ids", foodCtrl.RPCGetByIds)
	g.PATCH("/rpc/foods/update/:id", foodCtrl.RPCUpdateFoodByIdAPI)

	// Foods
	foods := g.Group("/foods")
	foodCtrl.SetupRoutes(foods)
}
