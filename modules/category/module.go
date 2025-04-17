package categorymodule

import (
	"github.com/gin-gonic/gin"
	categoryHttpgin "github.com/ntttrang/go-food-delivery-backend-service/modules/category/infras/controller/http-gin"
	categoryRepo "github.com/ntttrang/go-food-delivery-backend-service/modules/category/infras/repository"
	categoryService "github.com/ntttrang/go-food-delivery-backend-service/modules/category/service"
	shareinfras "github.com/ntttrang/go-food-delivery-backend-service/shared/infras"
)

func SetupCategoryModule(appCtx shareinfras.IAppContext, g *gin.RouterGroup) {
	dbCtx := appCtx.DbContext()

	// Setup controller
	catRepo := categoryRepo.NewCategoryRepo(dbCtx)
	bulkCreateCmdHdl := categoryService.NewBulkCreateCommandHandler(catRepo)
	createCmdHdl := categoryService.NewCreateCommandHandler(catRepo)
	listCmdHdl := categoryService.NewListCommandHandler(catRepo)
	getDetailCmdHdl := categoryService.NewGetDetailQueryHandler(catRepo)
	updateCmdHdl := categoryService.NewUpdateCommandHandler(catRepo)
	deleteCmdHdl := categoryService.NewDeleteByIdCommandHandler(catRepo)

	catCtl := categoryHttpgin.NewCategoryHttpController(bulkCreateCmdHdl, createCmdHdl, listCmdHdl, getDetailCmdHdl, updateCmdHdl, deleteCmdHdl, catRepo)

	// Setup router
	categories := g.Group("/categories")
	catCtl.SetupRoutes(categories)

	// Setup router ( internal)
	catCtl.SetupRoutesRPC(g)
}
