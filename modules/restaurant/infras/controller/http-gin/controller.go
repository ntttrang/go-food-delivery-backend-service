package httpgin

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ntttrang/go-food-delivery-backend-service/middleware"
	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	sharedrpc "github.com/ntttrang/go-food-delivery-backend-service/shared/infras/rpc"
)

type ICreateCommandHandler interface {
	Execute(ctx context.Context, data *restaurantmodel.RestaurantInsertDto) error
}

type IListQueryHandler interface {
	Execute(ctx context.Context, req restaurantmodel.RestaurantListReq) (restaurantmodel.RestaurantSearchRes, error)
}

type IGetDetailQueryHandler interface {
	Execute(ctx context.Context, req restaurantmodel.RestaurantDetailReq) (restaurantmodel.RestaurantDetailRes, error)
}

type IUpdateRestaurantCommandHandler interface {
	Execute(ctx context.Context, req restaurantmodel.RestaurantUpdateReq) error
}

type IDeleteCommandHandler interface {
	Execute(ctx context.Context, req restaurantmodel.RestaurantDeleteReq) error
}

type IAddFavoritesCommandHandler interface {
	Execute(ctx context.Context, req restaurantmodel.RestaurantLike) (*string, error)
}
type IListFavoritesQueryHandler interface {
	Execute(ctx context.Context, req restaurantmodel.FavoriteRestaurantListReq) (restaurantmodel.RestaurantSearchRes, error)
}

type ICreateRestaurantCommentCommandHandler interface {
	Execute(ctx context.Context, req *restaurantmodel.RestaurantCommentCreateReq) error
}

type IListRestaurantCommentsQueryHandler interface {
	Execute(ctx context.Context, req restaurantmodel.RestaurantRatingListReq) (*restaurantmodel.RestaurantRatingListRes, error)
}

type IDeleteCommentCommandHandler interface {
	Execute(ctx context.Context, req restaurantmodel.RestaurantDeleteCommentReq) error
}

type ICreateMenuItemCommandHandler interface {
	Execute(ctx context.Context, req *restaurantmodel.MenuItemCreateReq) error
}

type IListMenuItemQueryHandler interface {
	Execute(ctx context.Context, restaurantId uuid.UUID) (*restaurantmodel.MenuItemListRes, error)
}

type IDeleteMenuItemCommandHandler interface {
	Execute(ctx context.Context, req *restaurantmodel.MenuItemCreateReq) error
}

type RestaurantHttpController struct {
	createCmdHdl      ICreateCommandHandler
	listQueryHdl      IListQueryHandler
	getDetailQueryHdl IGetDetailQueryHandler
	updateCmdHdl      IUpdateRestaurantCommandHandler
	deleteCmdHdl      IDeleteCommandHandler

	addFavoritesCmdHdl         IAddFavoritesCommandHandler
	favoriteRestaurantQueryHdl IListFavoritesQueryHandler

	createCommentRestaurantCmdHandler ICreateRestaurantCommentCommandHandler
	listRestaurantCommentQueryHandler IListRestaurantCommentsQueryHandler
	deleteRestaurantCmdHdl            IDeleteCommentCommandHandler

	createMenuItemCmdHdl     ICreateMenuItemCommandHandler
	listMenuItemQueryHandler IListMenuItemQueryHandler
	deleteMenuItemCmdHdl     IDeleteMenuItemCommandHandler
}

func NewRestaurantHttpController(createCmdHdl ICreateCommandHandler, listQueryHdl IListQueryHandler, getDetailQueryHdl IGetDetailQueryHandler,
	updateCmdHdl IUpdateRestaurantCommandHandler, deleteCmdHdl IDeleteCommandHandler,
	addFavoritesCmdHdl IAddFavoritesCommandHandler, favoriteRestaurantQueryHdl IListFavoritesQueryHandler,
	createCommentRestaurantCmdHandler ICreateRestaurantCommentCommandHandler, listRestaurantCommentQueryHandler IListRestaurantCommentsQueryHandler, deleteRestaurantCmdHdl IDeleteCommentCommandHandler,
	createMenuItemCmdHdl ICreateMenuItemCommandHandler, listMenuItemQueryHandler IListMenuItemQueryHandler, deleteMenuItemCmdHdl IDeleteMenuItemCommandHandler) *RestaurantHttpController {
	return &RestaurantHttpController{
		createCmdHdl:      createCmdHdl,
		listQueryHdl:      listQueryHdl,
		getDetailQueryHdl: getDetailQueryHdl,
		updateCmdHdl:      updateCmdHdl,
		deleteCmdHdl:      deleteCmdHdl,

		addFavoritesCmdHdl:         addFavoritesCmdHdl,
		favoriteRestaurantQueryHdl: favoriteRestaurantQueryHdl,

		createCommentRestaurantCmdHandler: createCommentRestaurantCmdHandler,
		listRestaurantCommentQueryHandler: listRestaurantCommentQueryHandler,
		deleteRestaurantCmdHdl:            deleteRestaurantCmdHdl,

		createMenuItemCmdHdl:     createMenuItemCmdHdl,
		listMenuItemQueryHandler: listMenuItemQueryHandler,
		deleteMenuItemCmdHdl:     deleteMenuItemCmdHdl,
	}
}

func (ctrl *RestaurantHttpController) SetupRoutes(g *gin.RouterGroup) {
	introspectRpcClient := sharedrpc.NewIntrospectRpcClient(os.Getenv("USER_SERVICE_URL"))
	// Restaurant
	g.POST("", middleware.Auth(introspectRpcClient), ctrl.CreateRestaurantAPI)
	g.GET("", ctrl.ListRestaurantsAPI)         // Query params
	g.GET("/:id", ctrl.GetRestaurantDetailAPI) // Path Variables
	g.PATCH("/:id", ctrl.UpdateRestaurantByIdAPI)
	g.DELETE("/:id", ctrl.DeleteRestaurantByIdAPI)

	// Favorites Restaurant
	g.POST("/favorites", middleware.Auth(introspectRpcClient), ctrl.UpdateFavoritesRestaurantAPI)
	g.GET("/favorites", middleware.Auth(introspectRpcClient), ctrl.ListFavoriteRestaurantsAPI)

	// Restaurant Comments
	g.POST("/comments", middleware.Auth(introspectRpcClient), ctrl.CreateRestaurantCommentAPI)
	g.GET("/comments", ctrl.ListRestaurantCommentAPI)
	g.DELETE("/comments/:id", ctrl.DeleteRestaurantCommentAPI)

	// Menu item (restaurant-food)
	g.POST("/menu-item", ctrl.CreateMenuItemAPI)
	g.GET("/menu-item/:restaurantId", ctrl.ListMenuItemAPI)
	g.DELETE("/menu-item", ctrl.DeleteMenuItemAPI)

	// Save restaurant_foods ->  update foods
}
