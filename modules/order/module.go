package ordermodule

import (
	"github.com/gin-gonic/gin"
	orderHttpgin "github.com/ntttrang/go-food-delivery-backend-service/modules/order/infras/controller/http-gin"
	orderRepo "github.com/ntttrang/go-food-delivery-backend-service/modules/order/infras/repository/gorm-mysql"
	grpcclient "github.com/ntttrang/go-food-delivery-backend-service/modules/order/infras/repository/grpc-client"
	rpcclient "github.com/ntttrang/go-food-delivery-backend-service/modules/order/infras/repository/rpc-client"
	orderService "github.com/ntttrang/go-food-delivery-backend-service/modules/order/service"
	shareComponent "github.com/ntttrang/go-food-delivery-backend-service/shared/component"
	shareinfras "github.com/ntttrang/go-food-delivery-backend-service/shared/infras"
)

func SetupOrderModule(appCtx shareinfras.IAppContext, g *gin.RouterGroup) {
	dbCtx := appCtx.DbContext()
	config := appCtx.GetConfig()

	// Setup repository
	orderRepo := orderRepo.NewOrderRepo(dbCtx)
	// Setup RPC clients
	//foodRpcClientRepo := rpcclient.NewFoodRPCClient(appCtx.GetConfig().FoodServiceURL)
	restaurantRpcClientRepo := rpcclient.NewRestaurantRPCClient(appCtx.GetConfig().RestaurantServiceURL)
	cartRpcClientRepo := rpcclient.NewCartRPCClient(config.CartServiceURL)
	cardRpcClientRepo := rpcclient.NewCardRPCClient(appCtx.GetConfig().PaymentServiceURL)
	userRpcClientRepo := rpcclient.NewUserRPCClient(appCtx.GetConfig().UserServiceURL)
	emailSvc := shareComponent.NewEmailService(appCtx.GetConfig().EmailConfig)

	// GRPC
	foodGrpcClient := grpcclient.NewFoodGRPCClient(appCtx.GetConfig().GrpcServiceURL)

	// Setup service
	cartConversionService := orderService.NewCartToOrderConversionService(cartRpcClientRepo, foodGrpcClient, restaurantRpcClientRepo)
	paymentService := orderService.NewPaymentProcessingService(
		cardRpcClientRepo,
	)

	inventoryService := orderService.NewInventoryCheckingService(
		foodGrpcClient,
		restaurantRpcClientRepo,
	)

	notificationService := orderService.NewOrderNotificationService(
		orderRepo,
		userRpcClientRepo,
		restaurantRpcClientRepo,
		emailSvc,
	)

	// Create command handler with all services
	createCmdHdl := orderService.NewCreateCommandHandler(
		orderRepo,
		paymentService,
		inventoryService,
		notificationService,
	)
	createFromCartCmdHdl := orderService.NewCreateFromCartCommandHandler(
		createCmdHdl,
		cartConversionService,
		paymentService,
		inventoryService,
		notificationService,
		appCtx.MsgBroker(),
	)
	listQueryHdl := orderService.NewListQueryHandler(orderRepo)
	getDetailQueryHdl := orderService.NewGetDetailQueryHandler(orderRepo)
	updateOrderStateCmdHdl := orderService.NewOrderStateManagementService(orderRepo, notificationService, appCtx.MsgBroker())
	deleteCmdHdl := orderService.NewDeleteCommandHandler(orderRepo)

	// Setup controller with unified state management
	orderCtl := orderHttpgin.NewOrderHttpController(
		createCmdHdl,
		createFromCartCmdHdl,
		listQueryHdl,
		getDetailQueryHdl,
		updateOrderStateCmdHdl,
		deleteCmdHdl,
	)

	// Setup routes
	orders := g.Group("/orders")
	orderCtl.SetupRoutes(orders)
}
