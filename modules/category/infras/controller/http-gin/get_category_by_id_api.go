package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	categorymodel "github.com/ntttrang/go-food-delivery-backend-service/modules/category/model"
)

func (ctrl *CategoryHttpController) GetCategoryByIdAPI(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := ctrl.getDetailCmdHdl.Execute(c.Request.Context(), categorymodel.CategoryDetailReq{Id: id})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": category})
}
