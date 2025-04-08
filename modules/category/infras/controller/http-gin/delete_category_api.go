package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	categorymodel "github.com/ntttrang/go-food-delivery-backend-service/modules/category/model"
)

func (ctrl *CategoryHttpController) DeleteCategoryByIdAPI(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

	}

	if err := ctrl.deleteCmdHdl.Execute(c.Request.Context(), categorymodel.CategoryDeleteReq{Id: id}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{"id": id}})
}
