package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	categorymodel "github.com/ntttrang/go-food-delivery-backend-service/modules/category/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

func (ctrl *CategoryHttpController) ListCategoryAPI(c *gin.Context) {
	var dto categorymodel.ListCategoryReq
	if err := c.ShouldBindJSON(&dto); err != nil {
		panic(datatype.ErrBadRequest.WithError(err.Error()))
	}

	dto.Paging.Process()

	data, total, err := ctrl.listCmdHdl.Execute(c.Request.Context(), dto)
	if err != nil {
		panic(err)
	}

	paging := &dto.Paging
	paging.Total = total
	c.JSON(http.StatusOK, gin.H{"data": data, "paging": paging})
}
