package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	usermodel "github.com/ntttrang/go-food-delivery-backend-service/modules/user/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

func (ctrl *UserHttpController) CreateUserAddrAPI(c *gin.Context) {
	var req usermodel.CreateUserAddrReq

	if err := c.ShouldBindJSON(&req); err != nil {
		panic(datatype.ErrBadRequest.WithError(err.Error()))
	}

	requester := c.MustGet(datatype.KeyRequester).(datatype.Requester)
	req.UserId = requester.Subject()

	if err := ctrl.createAddrCmdHdl.Execute(c.Request.Context(), &req); err != nil {
		panic(err)
	}

	c.JSON(http.StatusCreated, gin.H{"data": req.Id})
}
