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
	rpcRestaurantRepo := rpcclient.NewRestaurantRPCClient(appCtx.GetConfig().RestaurantServiceURL)

	// Setup command handlers
	createCmdHdl := cartService.NewCreateCommandHandler(repo, rpcFoodRepo)
	listQueryHdl := cartService.NewListQueryHandler(repo, rpcRestaurantRepo)
	listCartItemQueryHdl := cartService.NewListCartItemQueryHandler(repo, rpcFoodRepo, rpcRestaurantRepo)
	getDetailQueryHdl := cartService.NewGetDetailQueryHandler(repo, rpcFoodRepo)
	updateCmdHdl := cartService.NewUpdateCommandHandler(repo)
	deleteCmdHdl := cartService.NewDeleteCommandHandler(repo)

	// Setup controller (RPC operations call repository directly, no service layer needed)
	cartCtl := cartHttpgin.NewCartHttpController(
		createCmdHdl,
		listQueryHdl,
		listCartItemQueryHdl,
		getDetailQueryHdl,
		updateCmdHdl,
		deleteCmdHdl,
		repo,
	)

	// Setup routes
	carts := g.Group("/carts")
	cartCtl.SetupRoutes(carts)
}
