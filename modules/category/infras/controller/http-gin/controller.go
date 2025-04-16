package httpgin

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	categorymodel "github.com/ntttrang/go-food-delivery-backend-service/modules/category/model"
)

type IBulkCreateCommandHandler interface {
	Execute(ctx context.Context, datas []categorymodel.CategoryInsertDto) ([]uuid.UUID, error)
}

type ICreateCommandHandler interface {
	Execute(ctx context.Context, data *categorymodel.CategoryInsertDto) error
}

type IListCommandHandler interface {
	Execute(ctx context.Context, req categorymodel.ListCategoryReq) ([]categorymodel.ListCategoryRes, int64, error)
}

type IGetDetailCommandHandler interface {
	Execute(ctx context.Context, req categorymodel.CategoryDetailReq) (categorymodel.CategoryDetailRes, error)
}

type IUpdateByIdCommandHandler interface {
	Execute(ctx context.Context, req categorymodel.CategoryUpdateReq) error
}

type IDeleteCommandHandler interface {
	Execute(ctx context.Context, req categorymodel.CategoryDeleteReq) error
}

type CategoryHttpController struct {
	bulkCreateCmdHdl         IBulkCreateCommandHandler
	createCmdHdl             ICreateCommandHandler
	listCmdHdl               IListCommandHandler
	getDetailCmdHdl          IGetDetailCommandHandler
	updateByIdCommandHandler IUpdateByIdCommandHandler
	deleteCmdHdl             IDeleteCommandHandler
}

func NewCategoryHttpController(bulkCreateCmdHdl IBulkCreateCommandHandler, createCmdHdl ICreateCommandHandler, listCmdHdl IListCommandHandler, getDetailCmdHdl IGetDetailCommandHandler,
	updateByIdCommandHandler IUpdateByIdCommandHandler, deleteCmdHdl IDeleteCommandHandler) *CategoryHttpController {
	return &CategoryHttpController{
		bulkCreateCmdHdl:         bulkCreateCmdHdl,
		createCmdHdl:             createCmdHdl,
		listCmdHdl:               listCmdHdl,
		getDetailCmdHdl:          getDetailCmdHdl,
		updateByIdCommandHandler: updateByIdCommandHandler,
		deleteCmdHdl:             deleteCmdHdl,
	}
}

func (ctrl *CategoryHttpController) SetupRoutes(g *gin.RouterGroup) {
	g.POST("", ctrl.CreateCategoryAPI)
	g.POST("bulk-insert", ctrl.CreateBulkCategoryAPI)
	g.GET("", ctrl.ListCategoryAPI)
	g.GET("/:id", ctrl.GetCategoryByIdAPI)
	g.PATCH("/:id", ctrl.UpdateCategoryByIdAPI)
	g.DELETE("/:id", ctrl.DeleteCategoryByIdAPI)
}
