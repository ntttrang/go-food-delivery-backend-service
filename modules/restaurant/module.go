package restaurantmodule

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	restaurantHttpgin "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/infras/controller/http-gin"
	elasticsearch "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/infras/repository/elasticsearch"
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

	// Setup basic handlers
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
	deleteMenuItemCmdHdl := restaurantService.NewDeleteMenuItemCommandHandler(restaurantFoodRepo)

	// Setup Elasticsearch if available
	var searchRestaurantQueryHandler *restaurantService.SearchRestaurantQueryHandler
	var syncRestaurantIndexCommandHandler *restaurantService.SyncRestaurantIndexCommandHandler

	// Try to initialize Elasticsearch client
	esClient, err := shareComponent.NewElasticsearchClient(appCtx.GetConfig().ElasticSearch)
	if err != nil {
		log.Printf("Elasticsearch initialization failed: %v. Search functionality will be disabled.", err)
	}

	// If Elasticsearch client was successfully created, set up search functionality
	if esClient != nil {
		// Create a new client with the restaurant index name
		restaurantIndexName := "restaurants"

		// Use the client with the restaurant index
		restaurantEsClient := esClient.WithIndex(restaurantIndexName)

		// Setup search repository
		restaurantSearchRepo := elasticsearch.NewRestaurantSearchRepo(restaurantEsClient)

		// Initialize the index with proper mapping
		if err := restaurantSearchRepo.Initialize(context.Background()); err != nil {
			log.Printf("Failed to initialize restaurant index: %v", err)
		}

		// Setup search handlers
		searchRestaurantQueryHandler = restaurantService.NewSearchRestaurantQueryHandler(restaurantSearchRepo)
		syncRestaurantIndexCommandHandler = restaurantService.NewSyncRestaurantIndexCommandHandler(restaurantRepo, restaurantSearchRepo)

		// Setup event handler for Elasticsearch operations
		// This handler would be used to hook into restaurant CRUD operations
		// to automatically update the Elasticsearch index
		_ = restaurantService.NewRestaurantElasticsearchHandler(restaurantSearchRepo)

		log.Println("Elasticsearch initialized successfully. Restaurant search functionality is enabled.")
	} else {
		log.Println("Elasticsearch client not available. Restaurant search functionality will be disabled.")
	}

	// Create dummy handlers if Elasticsearch is not available
	if searchRestaurantQueryHandler == nil {
		// Create a dummy search handler that returns empty results
		searchRestaurantQueryHandler = restaurantService.NewSearchRestaurantQueryHandler(nil)
	}

	if syncRestaurantIndexCommandHandler == nil {
		// Create a dummy sync handler with nil repositories
		syncRestaurantIndexCommandHandler = restaurantService.NewSyncRestaurantIndexCommandHandler(nil, nil)
	}

	resCtl := restaurantHttpgin.NewRestaurantHttpController(
		createCmdHdl, listQueryHdl, getDetailQueryHdl, updateCmdHdl, deleteCmdHdl,
		createRestaurantFavoriteCmdl, favoriteRestaurantQueryHdl,
		createCommentRestaurantCmdl, listCommentRestaurantCmdl, deleteCommentRestaurantCmdl,
		createMenuItemCmdHdl, listMenuItemCmdHdl, deleteMenuItemCmdHdl,
		searchRestaurantQueryHandler, syncRestaurantIndexCommandHandler,
	)

	restaurants := g.Group("/restaurants")
	resCtl.SetupRoutes(restaurants)
}
