package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

func (ctrl *RestaurantHttpController) ListRestaurantCommentCommandHandler(c *gin.Context) {
	var req restaurantmodel.RestaurantRatingListReq
	if err := c.ShouldBind(&req); err != nil {
		panic(datatype.ErrBadRequest.WithError(err.Error()))
	}
	var pagingDto = req.PagingDto
	pagingDto.Process()
	req.PagingDto = pagingDto

	result, err := ctrl.listRestaurantCommentQueryCmdHandler.Execute(c.Request.Context(), req)

	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}
