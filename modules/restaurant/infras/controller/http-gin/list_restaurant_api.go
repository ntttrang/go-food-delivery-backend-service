package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

func (ctrl *RestaurantHttpController) ListRestaurantsAPI(c *gin.Context) {
	var searchDto restaurantmodel.RestaurantSearchDto
	var pagingDto sharedModel.PagingDto
	var sortingDto sharedModel.SortingDto
	if err := c.ShouldBind(&pagingDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.ShouldBind(&searchDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.ShouldBind(&sortingDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pagingDto.Process()

	req := restaurantmodel.RestaurantListReq{
		RestaurantSearchDto: searchDto,
		PagingDto:           pagingDto,
		SortingDto:          sortingDto,
	}

	result, err := ctrl.listQueryHdl.Execute(c.Request.Context(), req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}
