package httpgin

import (
	"context"

	"github.com/gin-gonic/gin"
	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/internal/model"
)

type IRestaurantService interface {
	CreateRestaurant(ctx context.Context, data *restaurantmodel.RestaurantInsertDto) error
}

type RestaurantHttpController struct {
	restaurantService IRestaurantService
}

func NewRestaurantHttpController(restaurantService IRestaurantService) *RestaurantHttpController {
	return &RestaurantHttpController{
		restaurantService: restaurantService,
	}
}

func (ctrl *RestaurantHttpController) SetupRoutes(g *gin.RouterGroup) {
	g.POST("", ctrl.CreateRestaurantAPI)
	// g.POST("bulk-insert", ctl.CreateBulkCategoryAPI)
	// g.POST("list", ctl.ListCategoryAPI)
	// g.GET("/:id", ctl.GetCategoryByIdAPI)
}
