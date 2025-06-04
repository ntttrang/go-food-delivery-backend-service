package ordermodule

import (
	"github.com/gin-gonic/gin"
	orderHttpgin "github.com/ntttrang/go-food-delivery-backend-service/modules/order/infras/controller/http-gin"
	orderRepo "github.com/ntttrang/go-food-delivery-backend-service/modules/order/infras/repository/gorm-mysql"
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
	foodRpcClientRepo := rpcclient.NewFoodRPCClient(appCtx.GetConfig().FoodServiceURL)
	restaurantRpcClientRepo := rpcclient.NewRestaurantRPCClient(appCtx.GetConfig().RestaurantServiceURL)
	cartRpcClientRepo := rpcclient.NewCartRPCClient(config.CartServiceURL)
	cardRpcClientRepo := rpcclient.NewCardRPCClient(appCtx.GetConfig().PaymentServiceURL)
	userRpcClientRepo := rpcclient.NewUserRPCClient(appCtx.GetConfig().UserServiceURL)
	emailSvc := shareComponent.NewEmailService(appCtx.GetConfig().EmailConfig)

	// Setup service
	cartConversionService := orderService.NewCartToOrderConversionService(cartRpcClientRepo, foodRpcClientRepo, restaurantRpcClientRepo)
	paymentService := orderService.NewPaymentProcessingService(
		cardRpcClientRepo,
	)

	inventoryService := orderService.NewInventoryCheckingService(
		foodRpcClientRepo,
		restaurantRpcClientRepo,
	)

	notificationService := orderService.NewOrderNotificationService(
		orderRepo,
		userRpcClientRepo,
		restaurantRpcClientRepo,
		emailSvc,
		nil, // smsSvc - TODO: implement when SMS service is ready
		nil, // pushSvc - TODO: implement when push notification service is ready
	)

	// Create command handler with all services
	createCmdHdl := orderService.NewCreateCommandHandler(
		orderRepo,
		paymentService,
		inventoryService,
		notificationService,
	)

	// Create the cart-to-order handler with all services
	createFromCartCmdHdl := orderService.NewCreateFromCartCommandHandler(
		createCmdHdl,
		cartConversionService,
		paymentService,
		inventoryService,
		notificationService,
	)

	listQueryHdl := orderService.NewListQueryHandler(orderRepo)
	getDetailQueryHdl := orderService.NewGetDetailQueryHandler(orderRepo)
	updateCmdHdl := orderService.NewUpdateCommandHandler(orderRepo)
	deleteCmdHdl := orderService.NewDeleteCommandHandler(orderRepo)

	// Setup controller
	orderCtl := orderHttpgin.NewOrderHttpController(
		createCmdHdl,
		createFromCartCmdHdl,
		listQueryHdl,
		getDetailQueryHdl,
		updateCmdHdl,
		deleteCmdHdl,
	)

	// Setup routes
	orders := g.Group("/orders")
	orderCtl.SetupRoutes(orders)
}
