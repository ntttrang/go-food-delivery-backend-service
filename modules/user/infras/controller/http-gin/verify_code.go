package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

func (ctrl *UserHttpController) VerifyCodeAPI(c *gin.Context) {
	requester := c.MustGet(datatype.KeyRequester).(datatype.Requester)
	userId := requester.Subject()
	code := c.Param("code")

	isValid, err := ctrl.verifyCode.Execute(c.Request.Context(), userId, code)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, datatype.ResponseSuccess(isValid))
}
