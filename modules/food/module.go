package foodmodule

import (
	"github.com/gin-gonic/gin"
	foodHttpgin "github.com/ntttrang/go-food-delivery-backend-service/modules/food/infras/controller/http-gin"
	repo "github.com/ntttrang/go-food-delivery-backend-service/modules/food/infras/repository"
	rpcclient "github.com/ntttrang/go-food-delivery-backend-service/modules/food/infras/repository/rpc-client"
	foodService "github.com/ntttrang/go-food-delivery-backend-service/modules/food/service"
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

	foodCtrl := foodHttpgin.NewFoodHttpController(
		createCmdHdl, listCmdHdl, getDetailCmdHdl, updateCmdHdl, deleteCmdHdl, foodRepo,
		createFoodFavoriteCmdl, favoriteFoodQueryHdl,
		createCommentFoodCmdl, listCommentFoodCmdl, deleteCommentFoodCmdl,
	)

	// Setup router
	// RPC
	g.POST("/rpc/foods/find-by-ids", foodCtrl.RPCGetByIds)
	g.PATCH("/rpc/foods/update/:id", foodCtrl.RPCUpdateFoodByIdAPI)

	// Foods
	foods := g.Group("/foods")
	foodCtrl.SetupRoutes(foods)
}
