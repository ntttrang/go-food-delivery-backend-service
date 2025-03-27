package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	categorymodel "github.com/ntttrang/go-food-delivery-backend-service/modules/category/model"
)

func (ctl *CategoryHttpController) ListCategoryAPI(c *gin.Context) {
	var dto categorymodel.ListCategoryReq

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dto.Paging.Process()

	data, total, err := ctl.listCmdHdl.Execute(c.Request.Context(), dto)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	paging := &dto.Paging
	paging.Total = total
	c.JSON(http.StatusOK, gin.H{"data": data, "paging": paging})
}
