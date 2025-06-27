package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ntttrang/go-food-delivery-backend-service/modules/payment/service"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

func (ctrl *CardController) UpdateCardStatus(c *gin.Context) {
	// Parse the ID from the URL
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		panic(datatype.ErrBadRequest.WithError(err.Error()))
	}

	// Parse the request body
	var req struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(datatype.ErrBadRequest.WithError(err.Error()))
	}

	// Execute the command
	err = ctrl.updateCardStatusHandler.Execute(c.Request.Context(), &service.CardUpdateStatusDto{
		ID:     id,
		Status: req.Status,
	})
	if err != nil {
		panic(err)
	}

	// Return the response
	c.JSON(http.StatusCreated, gin.H{"data": id})
}
