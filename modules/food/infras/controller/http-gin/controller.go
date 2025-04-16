package httpgin

import (
	"context"

	"github.com/gin-gonic/gin"
	foodmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/food/model"
)

type ICreateCommandHandler interface {
	Execute(ctx context.Context, data *foodmodel.FoodInsertDto) error
}

type IListCommandHandler interface {
	Execute(ctx context.Context, req foodmodel.ListFoodReq) ([]foodmodel.ListFoodRes, int64, error)
}

type IGetDetailCommandHandler interface {
	Execute(ctx context.Context, req foodmodel.FoodDetailReq) (foodmodel.FoodDetailRes, error)
}

type IUpdateByIdCommandHandler interface {
	Execute(ctx context.Context, req foodmodel.FoodUpdateReq) error
}

type IDeleteCommandHandler interface {
	Execute(ctx context.Context, req foodmodel.FoodDeleteReq) error
}

type FoodHttpController struct {
	createCmdHdl             ICreateCommandHandler
	listCmdHdl               IListCommandHandler
	getDetailCmdHdl          IGetDetailCommandHandler
	updateByIdCommandHandler IUpdateByIdCommandHandler
	deleteCmdHdl             IDeleteCommandHandler
}

func NewFoodHttpController(createCmdHdl ICreateCommandHandler, listCmdHdl IListCommandHandler, getDetailCmdHdl IGetDetailCommandHandler,
	updateByIdCommandHandler IUpdateByIdCommandHandler, deleteCmdHdl IDeleteCommandHandler) *FoodHttpController {
	return &FoodHttpController{
		createCmdHdl:             createCmdHdl,
		listCmdHdl:               listCmdHdl,
		getDetailCmdHdl:          getDetailCmdHdl,
		updateByIdCommandHandler: updateByIdCommandHandler,
		deleteCmdHdl:             deleteCmdHdl,
	}
}

func (ctrl *FoodHttpController) SetupRoutes(g *gin.RouterGroup) {
	g.POST("", ctrl.CreateFoodAPI)
	g.GET("", ctrl.ListFoodAPI)
	g.GET("/:id", ctrl.GetFoodByIdAPI)
	g.PATCH("/:id", ctrl.UpdateFoodByIdAPI)
	g.DELETE("/:id", ctrl.DeleteFoodByIdAPI)
}
