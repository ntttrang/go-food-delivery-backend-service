package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

func (ctrl *UserHttpController) GenerateCodeAPI(c *gin.Context) {
	requester := c.MustGet(datatype.KeyRequester).(datatype.Requester)
	userId := requester.Subject()

	verifyCode, err := ctrl.generateCode.Execute(c.Request.Context(), userId)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, datatype.ResponseSuccess(verifyCode))
}
