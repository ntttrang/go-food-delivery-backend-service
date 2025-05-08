package httpgin

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ntttrang/go-food-delivery-backend-service/middleware"
	foodmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/food/model"
	"github.com/ntttrang/go-food-delivery-backend-service/modules/food/service"
	sharerpc "github.com/ntttrang/go-food-delivery-backend-service/shared/infras/rpc"
)

type ICreateCommandHandler interface {
	Execute(ctx context.Context, data *service.FoodInsertDto) error
}

type IListCommandHandler interface {
	Execute(ctx context.Context, req service.ListFoodReq) (*service.ListFoodRes, error)
}

type IGetDetailCommandHandler interface {
	Execute(ctx context.Context, req service.FoodDetailReq) (*service.FoodDetailRes, error)
}

type IUpdateByIdCommandHandler interface {
	Execute(ctx context.Context, req service.FoodUpdateReq) error
}

type IDeleteCommandHandler interface {
	Execute(ctx context.Context, req service.FoodDeleteReq) error
}

type IRepoRPCFood interface {
	FindByIds(ctx context.Context, ids []uuid.UUID) ([]foodmodel.FoodInfoDto, error)
}

type IAddFavoritesCommandHandler interface {
	Execute(ctx context.Context, req foodmodel.FoodLike) (*string, error)
}
type IListFavoritesQueryHandler interface {
	Execute(ctx context.Context, req service.FavoriteFoodListReq) (*service.ListFoodRes, error)
}

type ICreateFoodCommentCommandHandler interface {
	Execute(ctx context.Context, req *service.FoodCommentCreateReq) error
}

type IListFoodCommentsQueryHandler interface {
	Execute(ctx context.Context, req service.FoodRatingListReq) (*service.FoodRatingListRes, error)
}

type IDeleteCommentCommandHandler interface {
	Execute(ctx context.Context, req service.FoodDeleteCommentReq) error
}

type ISearchFoodQueryHandler interface {
	Execute(ctx context.Context, req service.FoodSearchReq) (*service.FoodSearchRes, error)
}

type ISyncFoodByIdCommandHandler interface {
	SyncFood(ctx context.Context, id uuid.UUID) error
}

type ISyncFoodIndexCommandHandler interface {
	SyncAll(ctx context.Context) error
}

type FoodHttpController struct {
	createCmdHdl             ICreateCommandHandler
	listCmdHdl               IListCommandHandler
	getDetailCmdHdl          IGetDetailCommandHandler
	updateByIdCommandHandler IUpdateByIdCommandHandler
	deleteCmdHdl             IDeleteCommandHandler
	rpcRepo                  IRepoRPCFood

	addFavoritesCmdHdl   IAddFavoritesCommandHandler
	favoriteFoodQueryHdl IListFavoritesQueryHandler

	createCommentFoodCmdHandler ICreateFoodCommentCommandHandler
	listFoodCommentQueryHandler IListFoodCommentsQueryHandler
	deleteFoodCmdHdl            IDeleteCommentCommandHandler

	searchFoodQueryHandler      ISearchFoodQueryHandler
	syncFoodByIdCommandHandler  ISyncFoodByIdCommandHandler
	syncFoodIndexCommandHandler ISyncFoodIndexCommandHandler
}

func NewFoodHttpController(createCmdHdl ICreateCommandHandler, listCmdHdl IListCommandHandler, getDetailCmdHdl IGetDetailCommandHandler,
	updateByIdCommandHandler IUpdateByIdCommandHandler, deleteCmdHdl IDeleteCommandHandler,
	rpcRepo IRepoRPCFood,
	addFavoritesCmdHdl IAddFavoritesCommandHandler, favoriteFoodQueryHdl IListFavoritesQueryHandler,
	createCommentFoodCmdHandler ICreateFoodCommentCommandHandler, listFoodCommentQueryHandler IListFoodCommentsQueryHandler, deleteFoodCmdHdl IDeleteCommentCommandHandler,
	searchFoodQueryHandler ISearchFoodQueryHandler, syncFoodByIdCommandHandler ISyncFoodByIdCommandHandler, syncFoodIndexCommandHandler ISyncFoodIndexCommandHandler) *FoodHttpController {
	return &FoodHttpController{
		createCmdHdl:             createCmdHdl,
		listCmdHdl:               listCmdHdl,
		getDetailCmdHdl:          getDetailCmdHdl,
		updateByIdCommandHandler: updateByIdCommandHandler,
		deleteCmdHdl:             deleteCmdHdl,
		rpcRepo:                  rpcRepo,

		addFavoritesCmdHdl:   addFavoritesCmdHdl,
		favoriteFoodQueryHdl: favoriteFoodQueryHdl,

		createCommentFoodCmdHandler: createCommentFoodCmdHandler,
		listFoodCommentQueryHandler: listFoodCommentQueryHandler,
		deleteFoodCmdHdl:            deleteFoodCmdHdl,

		searchFoodQueryHandler:      searchFoodQueryHandler,
		syncFoodByIdCommandHandler:  syncFoodByIdCommandHandler,
		syncFoodIndexCommandHandler: syncFoodIndexCommandHandler,
	}
}

func (ctrl *FoodHttpController) SetupRoutes(g *gin.RouterGroup) {
	g.POST("", ctrl.CreateFoodAPI)
	g.GET("", ctrl.ListFoodAPI)
	g.GET("/:id", ctrl.GetFoodByIdAPI)
	g.PATCH("/:id", ctrl.UpdateFoodByIdAPI)
	g.DELETE("/:id", ctrl.DeleteFoodByIdAPI)

	// Favorites Food
	introspectRpcClient := sharerpc.NewIntrospectRpcClient(os.Getenv("USER_SERVICE_URL"))
	g.POST("/favorites", middleware.Auth(introspectRpcClient), ctrl.UpdateFavoritesFoodAPI)
	g.GET("/favorites", middleware.Auth(introspectRpcClient), ctrl.ListFavoriteFoodAPI)

	// Food Comments
	g.POST("/comments", middleware.Auth(introspectRpcClient), ctrl.CreateFoodCommentAPI)
	g.GET("/comments", ctrl.ListFoodCommentAPI)
	g.DELETE("/comments/:id", ctrl.DeleteFoodCommentAPI)

	// Search
	g.POST("/search", ctrl.SearchFoodAPI)

	// Admin endpoints for Elasticsearch index management
	adminGroup := g.Group("/admin")
	adminGroup.POST("/sync-index", middleware.Auth(introspectRpcClient), ctrl.SyncFoodIndexAPI)
	adminGroup.POST("/sync-index/:id", middleware.Auth(introspectRpcClient), ctrl.SyncFoodByIdAPI)
}
