package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ntttrang/go-food-delivery-backend-service/modules/food/service"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

func (ctrl *FoodHttpController) ListFoodAPI(c *gin.Context) {
	var searchDto service.SearchFoodDto
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

	req := service.ListFoodReq{
		SearchFoodDto: searchDto,
		PagingDto:     pagingDto,
		SortingDto:    sortingDto,
	}

	data, err := ctrl.listCmdHdl.Execute(c.Request.Context(), req)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}
