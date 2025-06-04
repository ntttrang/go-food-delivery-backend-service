package httpgin

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ntttrang/go-food-delivery-backend-service/middleware"
	"github.com/ntttrang/go-food-delivery-backend-service/modules/order/service"
	sharedrpc "github.com/ntttrang/go-food-delivery-backend-service/shared/infras/rpc"
)

type ICreateCommandHandler interface {
	Execute(ctx context.Context, data *service.OrderCreateDto) (string, error)
}

type ICreateFromCartCommandHandler interface {
	ExecuteFromCart(ctx context.Context, data *service.OrderCreateFromCartDto) (string, error)
}

type IListQueryHandler interface {
	Execute(ctx context.Context, req service.OrderListReq) (*service.OrderListRes, error)
}

type IGetDetailQueryHandler interface {
	Execute(ctx context.Context, req service.OrderDetailReq) (*service.OrderDetailRes, error)
}

type IUpdateCommandHandler interface {
	Execute(ctx context.Context, req service.OrderUpdateReq) error
}

type IDeleteCommandHandler interface {
	Execute(ctx context.Context, req service.OrderDeleteReq) error
}

type OrderHttpController struct {
	createCmdHdl         ICreateCommandHandler
	createFromCartCmdHdl ICreateFromCartCommandHandler
	listQueryHdl         IListQueryHandler
	getDetailQueryHdl    IGetDetailQueryHandler
	updateCmdHdl         IUpdateCommandHandler
	deleteCmdHdl         IDeleteCommandHandler
}

func NewOrderHttpController(
	createCmdHdl ICreateCommandHandler,
	createFromCartCmdHdl ICreateFromCartCommandHandler,
	listQueryHdl IListQueryHandler,
	getDetailQueryHdl IGetDetailQueryHandler,
	updateCmdHdl IUpdateCommandHandler,
	deleteCmdHdl IDeleteCommandHandler,
) *OrderHttpController {
	return &OrderHttpController{
		createCmdHdl:         createCmdHdl,
		createFromCartCmdHdl: createFromCartCmdHdl,
		listQueryHdl:         listQueryHdl,
		getDetailQueryHdl:    getDetailQueryHdl,
		updateCmdHdl:         updateCmdHdl,
		deleteCmdHdl:         deleteCmdHdl,
	}
}

func (ctrl *OrderHttpController) SetupRoutes(g *gin.RouterGroup) {
	introspectRpcClient := sharedrpc.NewIntrospectRpcClient(os.Getenv("USER_SERVICE_URL"))

	// Order routes
	// NOTE: This CreateOrderAPI is restricted to ADMIN users only via API layer
	// Use CreateOrderFromCartAPI for standard customer order creation
	g.POST("", middleware.Auth(introspectRpcClient), ctrl.CreateOrderAPI)
	g.POST("/from-cart", middleware.Auth(introspectRpcClient), ctrl.CreateOrderFromCartAPI)
	g.GET("", ctrl.ListOrdersAPI)
	g.GET("/:id", ctrl.GetOrderDetailAPI)
	g.PATCH("/:id", ctrl.UpdateOrderAPI)
	g.DELETE("/:id", ctrl.DeleteOrderAPI)

	// Order state management routes
	g.PATCH("/:id/state", middleware.Auth(introspectRpcClient), ctrl.UpdateOrderStateAPI)
	g.PATCH("/:id/assign-shipper", middleware.Auth(introspectRpcClient), ctrl.AssignShipperAPI)
	g.PATCH("/:id/payment-status", middleware.Auth(introspectRpcClient), ctrl.UpdatePaymentStatusAPI)
}
