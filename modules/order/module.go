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

	// Setup repository
	repo := orderRepo.NewOrderRepo(dbCtx)

	// Setup command handlers
	createCmdHdl := orderService.NewCreateCommandHandler(repo)
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
