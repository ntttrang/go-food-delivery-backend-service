package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntttrang/go-food-delivery-backend-service/modules/category/service"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

func (ctrl *CategoryHttpController) ListCategoryAPI(c *gin.Context) {
	var searchDto service.SearchCategoryDto
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
	req := service.ListCategoryReq{
		SearchCategoryDto: searchDto,
		PagingDto:         pagingDto,
		SortingDto:        sortingDto,
	}

	data, total, err := ctrl.listCmdHdl.Execute(c.Request.Context(), req)
	if err != nil {
		panic(err)
	}
	if data == nil {
		data = make([]service.ListCategoryRes, 0)
	}

	paging := &pagingDto
	paging.Total = total
	c.JSON(http.StatusOK, gin.H{"data": data, "paging": paging})
}
