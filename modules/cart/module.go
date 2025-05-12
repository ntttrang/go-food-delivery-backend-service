package cartmodule

import (
	"github.com/gin-gonic/gin"
	cartHttpgin "github.com/ntttrang/go-food-delivery-backend-service/modules/cart/infras/controller/http-gin"
	cartRepo "github.com/ntttrang/go-food-delivery-backend-service/modules/cart/infras/repository/gorm-mysql"
	rpcclient "github.com/ntttrang/go-food-delivery-backend-service/modules/cart/infras/repository/rpcclient"
	cartService "github.com/ntttrang/go-food-delivery-backend-service/modules/cart/service"
	shareinfras "github.com/ntttrang/go-food-delivery-backend-service/shared/infras"
)

func SetupCartModule(appCtx shareinfras.IAppContext, g *gin.RouterGroup) {
	dbCtx := appCtx.DbContext()

	// Setup repository
	repo := cartRepo.NewCartRepo(dbCtx)
	rpcFoodRepo := rpcclient.NewFoodRPCClient(appCtx.GetConfig().FoodServiceURL)

	// Setup command handlers
	createCmdHdl := cartService.NewCreateCommandHandler(repo, rpcFoodRepo)
	listQueryHdl := cartService.NewListQueryHandler(repo, rpcFoodRepo)
	getDetailQueryHdl := cartService.NewGetDetailQueryHandler(repo, rpcFoodRepo)
	updateCmdHdl := cartService.NewUpdateCommandHandler(repo)
	deleteCmdHdl := cartService.NewDeleteCommandHandler(repo)

	// Setup controller
	cartCtl := cartHttpgin.NewCartHttpController(
		createCmdHdl,
		listQueryHdl,
		getDetailQueryHdl,
		updateCmdHdl,
		deleteCmdHdl,
	)

	// Setup routes
	carts := g.Group("/carts")
	cartCtl.SetupRoutes(carts)
}
