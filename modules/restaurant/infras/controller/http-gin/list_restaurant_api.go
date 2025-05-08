package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	restaurantservice "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/service"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

func (ctrl *RestaurantHttpController) ListRestaurantsAPI(c *gin.Context) {
	var searchDto restaurantservice.RestaurantSearchDto
	var pagingDto sharedModel.PagingDto
	var sortingDto sharedModel.SortingDto
	if err := c.ShouldBind(&pagingDto); err != nil {
		panic(datatype.ErrBadRequest.WithError(err.Error()))
	}

	if err := c.ShouldBind(&searchDto); err != nil {
		panic(datatype.ErrBadRequest.WithError(err.Error()))
	}

	if err := c.ShouldBind(&sortingDto); err != nil {
		panic(datatype.ErrBadRequest.WithError(err.Error()))
	}

	pagingDto.Process()

	req := restaurantservice.RestaurantListReq{
		RestaurantSearchDto: searchDto,
		PagingDto:           pagingDto,
		SortingDto:          sortingDto,
	}

	result, err := ctrl.listQueryHdl.Execute(c.Request.Context(), req)

	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}
