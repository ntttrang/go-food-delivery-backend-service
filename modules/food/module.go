package foodmodule

import (
	"log"

	"github.com/gin-gonic/gin"
	foodHttpgin "github.com/ntttrang/go-food-delivery-backend-service/modules/food/infras/controller/http-gin"
	repo "github.com/ntttrang/go-food-delivery-backend-service/modules/food/infras/repository"
	elasticsearch "github.com/ntttrang/go-food-delivery-backend-service/modules/food/infras/repository/elasticsearch"
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

	// Setup Elasticsearch if available
	var searchFoodQueryHandler *foodService.SearchFoodQueryHandler
	var syncFoodIndexCommandHandler *foodService.SyncFoodIndexCommandHandler

	// Try to initialize Elasticsearch client
	esClient, err := shareComponent.NewElasticsearchClient(appCtx.GetConfig().ElasticSearch)
	if err != nil {
		panic(err)
	}

	// If Elasticsearch client was successfully created, set up search functionality
	if esClient != nil {
		// Setup search repository
		foodSearchRepo := elasticsearch.NewFoodSearchRepo(esClient)

		// Setup search handlers
		searchFoodQueryHandler = foodService.NewSearchFoodQueryHandler(foodSearchRepo)
		syncFoodIndexCommandHandler = foodService.NewSyncFoodIndexCommandHandler(foodRepo, foodSearchRepo)
		log.Println("Elasticsearch initialized successfully. Search functionality is enabled.")
	} else {
		log.Println("Elasticsearch client not available. Search functionality will be disabled.")
	}

	// Create a dummy search handler if Elasticsearch is not available
	if searchFoodQueryHandler == nil {
		searchFoodQueryHandler = &foodService.SearchFoodQueryHandler{}
	}

	if syncFoodIndexCommandHandler == nil {
		syncFoodIndexCommandHandler = &foodService.SyncFoodIndexCommandHandler{}
	}

	foodCtrl := foodHttpgin.NewFoodHttpController(
		createCmdHdl, listCmdHdl, getDetailCmdHdl, updateCmdHdl, deleteCmdHdl, foodRepo,
		createFoodFavoriteCmdl, favoriteFoodQueryHdl,
		createCommentFoodCmdl, listCommentFoodCmdl, deleteCommentFoodCmdl,
		searchFoodQueryHandler, syncFoodIndexCommandHandler,
	)

	// Setup router
	// RPC
	g.POST("/rpc/foods/find-by-ids", foodCtrl.RPCGetByIds)
	g.PATCH("/rpc/foods/update/:id", foodCtrl.RPCUpdateFoodByIdAPI)

	// Foods
	foods := g.Group("/foods")
	foodCtrl.SetupRoutes(foods)
}
