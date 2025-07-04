package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	service "github.com/ntttrang/go-food-delivery-backend-service/modules/user/service"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

func (ctrl *UserHttpController) UpdateUseAPI(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		panic(datatype.ErrBadRequest.WithError(err.Error()))
	}

	var req service.UpdateUserReq

	if err := c.ShouldBindJSON(&req); err != nil {
		panic(datatype.ErrBadRequest.WithError(err.Error()))
	}

	req.Id = id

	// Get requester from context and pass it to the service
	requester := c.MustGet(datatype.KeyRequester).(datatype.Requester)
	req.Requester = requester

	if err := ctrl.updateCmdHdl.Execute(c.Request.Context(), req); err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"data": true})
}
