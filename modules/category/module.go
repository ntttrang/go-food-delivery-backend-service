package categorymodule

import (
	"github.com/gin-gonic/gin"
	categoryHttpgin "github.com/ntttrang/go-food-delivery-backend-service/modules/category/infras/controller/http-gin"
	categoryRepository "github.com/ntttrang/go-food-delivery-backend-service/modules/category/infras/repository"
	categoryService "github.com/ntttrang/go-food-delivery-backend-service/modules/category/internal/service"
	"gorm.io/gorm"
)

func SetupCategoryModule(db *gorm.DB, g *gin.RouterGroup) {
	catRepo := categoryRepository.NewCategoryRepository(db)
	catService := categoryService.NewCategoryService(catRepo)
	catCtl := categoryHttpgin.NewCategoryHttpController(catService)

	catCtl.SetupRoutes(g)
}
