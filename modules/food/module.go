package foodmodule

import (
	"github.com/gin-gonic/gin"
	foodHttpgin "github.com/ntttrang/go-food-delivery-backend-service/modules/food/infras/controller/http-gin"
	foodRepo "github.com/ntttrang/go-food-delivery-backend-service/modules/food/infras/repository"
	foodService "github.com/ntttrang/go-food-delivery-backend-service/modules/food/service"
	shareinfras "github.com/ntttrang/go-food-delivery-backend-service/shared/infras"
)

func SetupFoodModule(appCtx shareinfras.IAppContext, g *gin.RouterGroup) {
	dbCtx := appCtx.DbContext()

	// Setup controller
	repo := foodRepo.NewFoodRepo(dbCtx)
	createCmdHdl := foodService.NewCreateCommandHandler(repo)
	listCmdHdl := foodService.NewListCommandHandler(repo)
	getDetailCmdHdl := foodService.NewGetDetailQueryHandler(repo)
	updateCmdHdl := foodService.NewUpdateCommandHandler(repo)
	deleteCmdHdl := foodService.NewDeleteByIdCommandHandler(repo)

	foodCtrl := foodHttpgin.NewFoodHttpController(createCmdHdl, listCmdHdl, getDetailCmdHdl, updateCmdHdl, deleteCmdHdl, repo)

	// Setup router
	// RPC
	g.POST("/rpc/foods/find-by-ids", foodCtrl.RPCGetByIds)
	g.PATCH("/rpc/foods/update/:id", foodCtrl.RPCUpdateFoodByIdAPI)

	// Foods
	foods := g.Group("/foods")
	foodCtrl.SetupRoutes(foods)
}
