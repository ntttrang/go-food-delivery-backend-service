package ordermodule

import (
	"github.com/gin-gonic/gin"
	orderHttpgin "github.com/ntttrang/go-food-delivery-backend-service/modules/order/infras/controller/http-gin"
	orderRepo "github.com/ntttrang/go-food-delivery-backend-service/modules/order/infras/repository/gorm-mysql"
	orderService "github.com/ntttrang/go-food-delivery-backend-service/modules/order/service"
	shareinfras "github.com/ntttrang/go-food-delivery-backend-service/shared/infras"
)

func SetupOrderModule(appCtx shareinfras.IAppContext, g *gin.RouterGroup) {
	dbCtx := appCtx.DbContext()
	// config := appCtx.GetConfig() // TODO: Use when implementing RPC integration

	// Setup repository
	repo := orderRepo.NewOrderRepo(dbCtx)

	// Setup RPC clients
	// cartRPCClient := sharerpc.NewCartRPCClient(config.CartServiceURL)

	// Setup services for order creation from cart
	// cartConversionService := orderService.NewCartToOrderConversionServiceWithRPC(cartRPCClient)
	// TODO: Integrate cart conversion service with full order creation flow

	// For now, use simple command handler until we implement full services
	createCmdHdl := orderService.NewCreateCommandHandlerSimple(repo)
	listQueryHdl := orderService.NewListQueryHandler(repo)
	getDetailQueryHdl := orderService.NewGetDetailQueryHandler(repo)
	updateCmdHdl := orderService.NewUpdateCommandHandler(repo)
	deleteCmdHdl := orderService.NewDeleteCommandHandler(repo)

	// Setup controller
	orderCtl := orderHttpgin.NewOrderHttpController(
		createCmdHdl,
		listQueryHdl,
		getDetailQueryHdl,
		updateCmdHdl,
		deleteCmdHdl,
	)

	// Setup routes
	orders := g.Group("/orders")
	orderCtl.SetupRoutes(orders)
}
