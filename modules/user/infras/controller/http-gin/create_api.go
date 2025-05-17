package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	service "github.com/ntttrang/go-food-delivery-backend-service/modules/user/service"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

func (ctrl *UserHttpController) CreateUserAPI(c *gin.Context) {
	var requestBodyData *service.CreateUserReq

	if err := c.ShouldBindJSON(&requestBodyData); err != nil {
		panic(datatype.ErrBadRequest.WithError(err.Error()))
	}

	// Get requester from context and pass it to the service
	requester := c.MustGet(datatype.KeyRequester).(datatype.Requester)
	requestBodyData.Requester = requester

	if err := ctrl.createCmdHdl.Execute(c.Request.Context(), requestBodyData); err != nil {
		panic(err)
	}

	c.JSON(http.StatusCreated, gin.H{"data": requestBodyData.Id})
}
