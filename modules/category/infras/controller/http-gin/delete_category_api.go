package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ntttrang/go-food-delivery-backend-service/modules/category/service"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

func (ctrl *CategoryHttpController) DeleteCategoryByIdAPI(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		panic(datatype.ErrBadRequest.WithError(err.Error()))
	}

	// Get requester from context for authorization
	requester := c.MustGet(datatype.KeyRequester).(datatype.Requester)

	if err := ctrl.deleteCmdHdl.Execute(c.Request.Context(), service.CategoryDeleteReq{Id: id, Requester: requester}); err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{"id": id}})
}
