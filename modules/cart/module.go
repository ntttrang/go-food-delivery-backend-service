package cartmodule

import (
	"github.com/gin-gonic/gin"
	cartHttpgin "github.com/ntttrang/go-food-delivery-backend-service/modules/cart/infras/controller/http-gin"
	cartRepo "github.com/ntttrang/go-food-delivery-backend-service/modules/cart/infras/repository/gorm-mysql"
	grpcclient "github.com/ntttrang/go-food-delivery-backend-service/modules/cart/infras/repository/grpc-client"
	rpcclient "github.com/ntttrang/go-food-delivery-backend-service/modules/cart/infras/repository/rpcclient"
	cartService "github.com/ntttrang/go-food-delivery-backend-service/modules/cart/service"
	shareinfras "github.com/ntttrang/go-food-delivery-backend-service/shared/infras"
)

func SetupCartModule(appCtx shareinfras.IAppContext, g *gin.RouterGroup) {
	dbCtx := appCtx.DbContext()

	// Setup repository
	repo := cartRepo.NewCartRepo(dbCtx)

	// RPC
	//rpcFoodRepo := rpcclient.NewFoodRPCClient(appCtx.GetConfig().FoodServiceURL)
	rpcRestaurantRepo := rpcclient.NewRestaurantRPCClient(appCtx.GetConfig().RestaurantServiceURL)

	// GRPC
	foodGrpcClient := grpcclient.NewFoodGRPCClient(appCtx.GetConfig().GrpcFoodServiceURL)

	// Setup command handlers
	createCmdHdl := cartService.NewCreateCommandHandler(repo, foodGrpcClient)
	listQueryHdl := cartService.NewListQueryHandler(repo, rpcRestaurantRepo)
	listCartItemQueryHdl := cartService.NewListCartItemQueryHandler(repo, foodGrpcClient, rpcRestaurantRepo)
	getDetailQueryHdl := cartService.NewGetDetailQueryHandler(repo, foodGrpcClient)
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

	// RPC endpoints for order service integration
	g.PATCH("/rpc/carts/update-status", cartCtl.UpdateCartStatusAPI)
	g.GET("/rpc/carts/cart-summary", cartCtl.GetCartSummaryAPI) // GET /cart-summaryrts?cardId=zzz?userId=xxx

	// Setup routes
	carts := g.Group("/carts")
	cartCtl.SetupRoutes(carts)
}
