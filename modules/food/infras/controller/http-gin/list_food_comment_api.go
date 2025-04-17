package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	foodmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/food/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

func (ctrl *FoodHttpController) ListFoodCommentAPI(c *gin.Context) {
	var req foodmodel.FoodRatingListReq
	if err := c.ShouldBind(&req); err != nil {
		panic(datatype.ErrBadRequest.WithError(err.Error()))
	}
	var pagingDto = req.PagingDto
	pagingDto.Process()
	req.PagingDto = pagingDto

	result, err := ctrl.listFoodCommentQueryHandler.Execute(c.Request.Context(), req)

	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}
