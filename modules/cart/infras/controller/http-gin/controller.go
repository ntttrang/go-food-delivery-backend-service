package httpgin

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ntttrang/go-food-delivery-backend-service/middleware"
	"github.com/ntttrang/go-food-delivery-backend-service/modules/cart/service"
	sharedrpc "github.com/ntttrang/go-food-delivery-backend-service/shared/infras/rpc"
)

type ICreateCommandHandler interface {
	Execute(ctx context.Context, data *service.CartUpsertDto) error
}

type IListCartQueryHandler interface {
	Execute(ctx context.Context, req service.CartListReq) (*service.CartListRes, error)
}

type IListCartItemQueryHandler interface {
	Execute(ctx context.Context, req service.CartItemListReq) (*service.CartItemListRes, error)
}

type IGetDetailQueryHandler interface {
	Execute(ctx context.Context, req service.CartDetailReq) (*service.CartDetailRes, error)
}

type IUpdateCommandHandler interface {
	Execute(ctx context.Context, req service.CartUpdateReq) error
}

type IDeleteCommandHandler interface {
	Execute(ctx context.Context, req service.CartDeleteReq) error
}

type CartHttpController struct {
	createCmdHdl         ICreateCommandHandler
	listCartQueryHdl     IListCartQueryHandler
	listCartItemQueryHdl IListCartItemQueryHandler
	getDetailQueryHdl    IGetDetailQueryHandler
	updateCmdHdl         IUpdateCommandHandler
	deleteCmdHdl         IDeleteCommandHandler
}

func NewCartHttpController(
	createCmdHdl ICreateCommandHandler,
	listCartQueryHdl IListCartQueryHandler,
	listCartItemQueryHdl IListCartItemQueryHandler,
	getDetailQueryHdl IGetDetailQueryHandler,
	updateCmdHdl IUpdateCommandHandler,
	deleteCmdHdl IDeleteCommandHandler,
) *CartHttpController {
	return &CartHttpController{
		createCmdHdl:         createCmdHdl,
		listCartQueryHdl:     listCartQueryHdl,
		listCartItemQueryHdl: listCartItemQueryHdl,
		getDetailQueryHdl:    getDetailQueryHdl,
		updateCmdHdl:         updateCmdHdl,
		deleteCmdHdl:         deleteCmdHdl,
	}
}

func (ctrl *CartHttpController) SetupRoutes(g *gin.RouterGroup) {
	introspectRpcClient := sharedrpc.NewIntrospectRpcClient(os.Getenv("USER_SERVICE_URL"))

	// Cart routes
	g.POST("", middleware.Auth(introspectRpcClient), ctrl.UpsertCartAPI)
	g.GET("", ctrl.ListCartAPI)
	g.GET("/cart-item", ctrl.ListCartItemAPI)
	g.GET("/:userId/:foodId", ctrl.GetCartByUserIdAndFoodIdAPI)
	g.PATCH("/:id", ctrl.UpdateCartByIdAPI)
	g.DELETE("/:userId/:foodId", ctrl.DeleteCartByIdAPI)
}
