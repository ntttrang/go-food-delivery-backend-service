package httpgin

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ntttrang/go-food-delivery-backend-service/middleware"
	gormmysql "github.com/ntttrang/go-food-delivery-backend-service/modules/cart/infras/repository/gorm-mysql"
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

// Repository interface for direct cart operations
type ICartRepository interface {
	UpdateCartStatusByCartID(ctx context.Context, cartID uuid.UUID, status string) error
	//FindCartItemsByCartID(ctx context.Context, cartID uuid.UUID, userID uuid.UUID) ([]CartItem, error)
	//FindCartByCartIDAndUserID(ctx context.Context, cartID uuid.UUID, userID uuid.UUID) (*CartItem, error)
	GetCartSummaryByCartID(ctx context.Context, cartID uuid.UUID, userID uuid.UUID) (*gormmysql.CartSummaryData, error)
}

type CartHttpController struct {
	createCmdHdl         ICreateCommandHandler
	listCartQueryHdl     IListCartQueryHandler
	listCartItemQueryHdl IListCartItemQueryHandler
	getDetailQueryHdl    IGetDetailQueryHandler
	updateCmdHdl         IUpdateCommandHandler
	deleteCmdHdl         IDeleteCommandHandler
	repo                 ICartRepository // Direct repository access for simple operations
}

func NewCartHttpController(
	createCmdHdl ICreateCommandHandler,
	listCartQueryHdl IListCartQueryHandler,
	listCartItemQueryHdl IListCartItemQueryHandler,
	getDetailQueryHdl IGetDetailQueryHandler,
	updateCmdHdl IUpdateCommandHandler,
	deleteCmdHdl IDeleteCommandHandler,
	repo ICartRepository,
) *CartHttpController {
	return &CartHttpController{
		createCmdHdl:         createCmdHdl,
		listCartQueryHdl:     listCartQueryHdl,
		listCartItemQueryHdl: listCartItemQueryHdl,
		getDetailQueryHdl:    getDetailQueryHdl,
		updateCmdHdl:         updateCmdHdl,
		deleteCmdHdl:         deleteCmdHdl,
		repo:                 repo,
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

	// RPC endpoints for order service integration
	g.PATCH("/update-status", ctrl.UpdateCartStatusAPI)
	g.GET("/cart-summary", ctrl.GetCartSummaryAPI) // GET /cart-summaryrts?cardId=zzz?userId=xxx
}
