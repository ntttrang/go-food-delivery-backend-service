package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	service "github.com/ntttrang/go-food-delivery-backend-service/modules/user/service"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

func (ctrl *UserHttpController) IntrospectTokenRpcAPI(c *gin.Context) {
	var req service.IntrospectReq
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(datatype.ErrBadRequest.WithWrap(err).WithDebug(err.Error()))
	}

	user, err := ctrl.introspectCmdHdl.Execute(c.Request.Context(), req)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, datatype.ResponseSuccess(user))

}
