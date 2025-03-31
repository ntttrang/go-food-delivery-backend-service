package httpgin

import (
	"context"

	"github.com/gin-gonic/gin"
	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
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
	Execute(ctx context.Context, req restaurantmodel.RestaurantLike) error
}

type IDeleteFavoritesCommandHandler interface {
	Execute(ctx context.Context, req restaurantmodel.RestaurantLike) error
}

type ICreateRestaurantCommentCommandHandler interface {
	Execute(ctx context.Context, req restaurantmodel.RestaurantCommentCreateReq) error
}

type IListRestaurantCommentsQueryHandler interface {
	Execute(ctx context.Context, req restaurantmodel.RestaurantRatingListReq) ([]restaurantmodel.RestaurantRatingListRes, error)
}

type RestaurantHttpController struct {
	createCmdHdl      ICreateCommandHandler
	listQueryHdl      IListQueryHandler
	getDetailQueryHdl IGetDetailQueryHandler
	updateCmdHdl      IUpdateRestaurantCommandHandler
	deleteCmdHdl      IDeleteCommandHandler

	addFavoritesCmdHdl    IAddFavoritesCommandHandler
	deleteFavoritesCmdHdl IDeleteFavoritesCommandHandler

	createCommentRestaurantCmdHandler    ICreateRestaurantCommentCommandHandler
	listRestaurantCommentQueryCmdHandler IListRestaurantCommentsQueryHandler
}

func NewRestaurantHttpController(createCmdHdl ICreateCommandHandler, listQueryHdl IListQueryHandler, getDetailQueryHdl IGetDetailQueryHandler,
	updateCmdHdl IUpdateRestaurantCommandHandler, deleteCmdHdl IDeleteCommandHandler,
	addFavoritesCmdHdl IAddFavoritesCommandHandler, deleteFavoritesCmdHdl IDeleteFavoritesCommandHandler,
	createCommentRestaurantCmdHandler ICreateRestaurantCommentCommandHandler, listRestaurantCommentQueryCmdHandler IListRestaurantCommentsQueryHandler) *RestaurantHttpController {
	return &RestaurantHttpController{
		createCmdHdl:      createCmdHdl,
		listQueryHdl:      listQueryHdl,
		getDetailQueryHdl: getDetailQueryHdl,
		updateCmdHdl:      updateCmdHdl,
		deleteCmdHdl:      deleteCmdHdl,

		addFavoritesCmdHdl:    addFavoritesCmdHdl,
		deleteFavoritesCmdHdl: deleteFavoritesCmdHdl,

		createCommentRestaurantCmdHandler:    createCommentRestaurantCmdHandler,
		listRestaurantCommentQueryCmdHandler: listRestaurantCommentQueryCmdHandler,
	}
}

func (ctrl *RestaurantHttpController) SetupRoutes(g *gin.RouterGroup) {
	g.POST("", ctrl.CreateRestaurantAPI)
	g.GET("list", ctrl.ListRestaurantsAPI)     // Query params
	g.GET("/:id", ctrl.GetRestaurantDetailAPI) // Path Variables
	g.PATCH("/:id", ctrl.UpdateRestaurantByIdAPI)
	g.DELETE("/:id", ctrl.DeleteRestaurantByIdAPI)

	g.POST("/favorites/add", ctrl.CreateFavoritesRestaurantAPI)
	g.DELETE("/favorites/delete", ctrl.DeleteFavoritesRestaurantAPI)

	g.POST("/comments/add", ctrl.CreateRestaurantCommentCommandHandler)
	g.GET("/comments/list", ctrl.ListRestaurantCommentCommandHandler)
	g.DELETE("/comments/delete", ctrl.DeleteFavoritesRestaurantAPI)
}
