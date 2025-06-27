package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntttrang/go-food-delivery-backend-service/modules/payment/service"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

func (ctrl *CardController) CreateCard(c *gin.Context) {
	// Get the requester from the context
	requester := c.MustGet(datatype.KeyRequester).(datatype.Requester)
	userID := requester.Subject()

	// Parse the request body
	var req service.CreateCardReq
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(datatype.ErrBadRequest.WithError(err.Error()))
	}

	// Set the user ID from the requester
	req.UserID = userID

	// Execute the command
	res, err := ctrl.createCardHandler.Execute(c.Request.Context(), &req)
	if err != nil {
		panic(err)
	}

	// Return the response
	c.JSON(http.StatusCreated, gin.H{"data": res.ID})
}
