package categorymodule

import (
	"github.com/gin-gonic/gin"
	categoryHttpgin "github.com/ntttrang/go-food-delivery-backend-service/modules/category/infras/controller/http-gin"
	categoryRepository "github.com/ntttrang/go-food-delivery-backend-service/modules/category/infras/repository"
	categoryService "github.com/ntttrang/go-food-delivery-backend-service/modules/category/service"
	"gorm.io/gorm"
)

func SetupCategoryModule(db *gorm.DB, g *gin.RouterGroup) {
	// Setup controller
	catRepo := categoryRepository.NewCategoryRepository(db)
	bulkCreateCmdHdl := categoryService.NewBulkCreateCommandHandler(catRepo)
	createCmdHdl := categoryService.NewCreateCommandHandler(catRepo)
	listCmdHdl := categoryService.NewListCommandHandler(catRepo)
	getDetailCmdHdl := categoryService.NewGetDetailQueryHandler(catRepo)
	updateCmdHdl := categoryService.NewUpdateCommandHandler(catRepo)
	deleteCmdHdl := categoryService.NewDeleteByIdCommandHandler(catRepo)
	catCtl := categoryHttpgin.NewCategoryHttpController(bulkCreateCmdHdl, createCmdHdl, listCmdHdl, getDetailCmdHdl, updateCmdHdl, deleteCmdHdl)

	// Setup router
	categories := g.Group("/categories")
	catCtl.SetupRoutes(categories)

	// Setup router ( internal)
	//catCtl.SetupRoutesRPC(g)
}
