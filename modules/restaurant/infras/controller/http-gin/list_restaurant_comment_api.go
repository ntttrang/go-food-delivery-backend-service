package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	restaurantservice "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/service"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

// ListRestaurantCommentAPI: List all comments belong to restaurant or user input
func (ctrl *RestaurantHttpController) ListRestaurantCommentAPI(c *gin.Context) {
	var req restaurantservice.RestaurantRatingListReq
	if err := c.ShouldBind(&req); err != nil {
		panic(datatype.ErrBadRequest.WithError(err.Error()))
	}
	var pagingDto = req.PagingDto
	pagingDto.Process()
	req.PagingDto = pagingDto

	result, err := ctrl.listRestaurantCommentQueryHandler.Execute(c.Request.Context(), req)

	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}
