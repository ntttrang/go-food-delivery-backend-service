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

type RestaurantHttpController struct {
	createCmdHdl      ICreateCommandHandler
	listQueryHdl      IListQueryHandler
	getDetailQueryHdl IGetDetailQueryHandler
	updateCmdHdl      IUpdateRestaurantCommandHandler
	deleteCmdHdl      IDeleteCommandHandler

	addFavoritesCmdHdl    IAddFavoritesCommandHandler
	deleteFavoritesCmdHdl IDeleteFavoritesCommandHandler
}

func NewRestaurantHttpController(createCmdHdl ICreateCommandHandler, listQueryHdl IListQueryHandler, getDetailQueryHdl IGetDetailQueryHandler,
	updateCmdHdl IUpdateRestaurantCommandHandler, deleteCmdHdl IDeleteCommandHandler,
	addFavoritesCmdHdl IAddFavoritesCommandHandler, deleteFavoritesCmdHdl IDeleteFavoritesCommandHandler) *RestaurantHttpController {
	return &RestaurantHttpController{
		createCmdHdl:      createCmdHdl,
		listQueryHdl:      listQueryHdl,
		getDetailQueryHdl: getDetailQueryHdl,
		updateCmdHdl:      updateCmdHdl,
		deleteCmdHdl:      deleteCmdHdl,

		addFavoritesCmdHdl:    addFavoritesCmdHdl,
		deleteFavoritesCmdHdl: deleteFavoritesCmdHdl,
	}
}

func (ctrl *RestaurantHttpController) SetupRoutes(g *gin.RouterGroup) {
	g.POST("", ctrl.CreateRestaurantAPI)
	g.GET("list", ctrl.ListRestaurantsAPI)     // Query params
	g.GET("/:id", ctrl.GetRestaurantDetailAPI) // Path Variables
	g.PATCH("/:id", ctrl.UpdateRestaurantByIdAPI)
	g.DELETE("/:id", ctrl.DeleteRestaurantByIdAPI)

	g.POST("/favorites/add", ctrl.AddFavoritesRestaurantAPI)
	g.DELETE("/favorites/delete", ctrl.DeleteFavoritesRestaurantAPI)
}
