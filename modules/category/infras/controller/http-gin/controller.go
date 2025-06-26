package httpgin

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	categorymodel "github.com/ntttrang/go-food-delivery-backend-service/modules/category/model"
	"github.com/ntttrang/go-food-delivery-backend-service/modules/category/service"
	sharedinfras "github.com/ntttrang/go-food-delivery-backend-service/shared/infras"
)

type ICreateCommandHandler interface {
	Execute(ctx context.Context, data *service.CategoryInsertDto) error
}

type IListCommandHandler interface {
	Execute(ctx context.Context, req service.ListCategoryReq) ([]service.ListCategoryRes, int64, error)
}

type IGetDetailCommandHandler interface {
	Execute(ctx context.Context, req service.CategoryDetailReq) (*service.CategoryDetailRes, error)
}

type IUpdateByIdCommandHandler interface {
	Execute(ctx context.Context, req service.CategoryUpdateReq) error
}

type IDeleteCommandHandler interface {
	Execute(ctx context.Context, req service.CategoryDeleteReq) error
}

type IRepoRPCCategory interface {
	FindByIds(ctx context.Context, ids []uuid.UUID) ([]categorymodel.Category, error)
}

type CategoryHttpController struct {
	createCmdHdl             ICreateCommandHandler
	listCmdHdl               IListCommandHandler
	getDetailCmdHdl          IGetDetailCommandHandler
	updateByIdCommandHandler IUpdateByIdCommandHandler
	deleteCmdHdl             IDeleteCommandHandler
	repoRPCCategory          IRepoRPCCategory
}

func NewCategoryHttpController(createCmdHdl ICreateCommandHandler, listCmdHdl IListCommandHandler, getDetailCmdHdl IGetDetailCommandHandler,
	updateByIdCommandHandler IUpdateByIdCommandHandler, deleteCmdHdl IDeleteCommandHandler,
	repoRPCCategory IRepoRPCCategory) *CategoryHttpController {
	return &CategoryHttpController{
		createCmdHdl:             createCmdHdl,
		listCmdHdl:               listCmdHdl,
		getDetailCmdHdl:          getDetailCmdHdl,
		updateByIdCommandHandler: updateByIdCommandHandler,
		deleteCmdHdl:             deleteCmdHdl,

		repoRPCCategory: repoRPCCategory,
	}
}

func (ctrl *CategoryHttpController) SetupRoutes(g *gin.RouterGroup, mldProvider sharedinfras.IMiddlewareProvider) {
	g.POST("", mldProvider.Auth(), ctrl.CreateCategoryAPI)
	g.GET("", ctrl.ListCategoryAPI)
	g.GET("/:id", ctrl.GetCategoryByIdAPI)
	g.PATCH("/:id", mldProvider.Auth(), ctrl.UpdateCategoryByIdAPI)
	g.DELETE("/:id", mldProvider.Auth(), ctrl.DeleteCategoryByIdAPI)
}

func (ctrl *CategoryHttpController) SetupRoutesRPC(g *gin.RouterGroup) {
	g.POST("/rpc/categories/find-by-ids", ctrl.RPCGetByIds)
}
