package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

type CallbackReq struct {
	Code  string `json:"code" form:"code"`
	State string `json:"state" form:"state"`
}

func (ctrl *UserHttpController) CallbackAPI(c *gin.Context) {
	var req CallbackReq
	if err := c.BindQuery(&req); err != nil {
		panic(datatype.ErrBadRequest.WithError("state is required"))
	}

	authRes, err := ctrl.signUpGgCmdHdl.AuthenticateByGoogle(c.Request.Context(), req.State, req.Code)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, datatype.ResponseSuccess(authRes))
}
