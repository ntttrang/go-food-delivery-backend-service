package httpgin

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	categorymodel "github.com/ntttrang/go-food-delivery-backend-service/modules/category/internal/model"
)

type ICategoryService interface {
	CreateCategory(ctx context.Context, data *categorymodel.CategoryInsertDto) error
	BulkInsert(ctx context.Context, datas []categorymodel.CategoryInsertDto) ([]uuid.UUID, error)
	ListCategories(ctx context.Context, req categorymodel.ListCategoryReq) ([]categorymodel.ListCategoryRes, int64, error)
}

type CategoryHttpController struct {
	catService ICategoryService
}

func NewCategoryHttpController(catService ICategoryService) *CategoryHttpController {
	return &CategoryHttpController{
		catService: catService,
	}
}

func (ctl *CategoryHttpController) SetupRoutes(g *gin.RouterGroup) {
	g.POST("", ctl.CreateCategoryAPI)
	g.POST("bulk-insert", ctl.CreateBulkCategoryAPI)
	g.POST("list", ctl.ListCategoryAPI)
	// g.GET("/:id", ctl.GetCategoryByIdAPI)
}
