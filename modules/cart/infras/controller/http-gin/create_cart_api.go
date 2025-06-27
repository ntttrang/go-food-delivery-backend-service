package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntttrang/go-food-delivery-backend-service/modules/cart/service"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

func (ctrl *CartHttpController) UpsertCartAPI(c *gin.Context) {
	var req service.CartUpsertDto

	if err := c.ShouldBindJSON(&req); err != nil {
		panic(datatype.ErrBadRequest.WithError(err.Error()))
	}

	// Get user ID from requester context
	requester := c.MustGet(datatype.KeyRequester).(datatype.Requester)
	req.UserID = requester.Subject()

	// Call business logic in service
	if err := ctrl.createCmdHdl.Execute(c.Request.Context(), &req); err != nil {
		panic(err)
	}

	c.JSON(http.StatusCreated, gin.H{"data": req.ID})
}
