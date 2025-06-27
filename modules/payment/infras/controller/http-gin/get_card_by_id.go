package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

func (ctrl *CardController) GetCardByID(c *gin.Context) {
	// Parse the ID from the URL
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		panic(datatype.ErrBadRequest.WithError(err.Error()))
	}

	// Execute the query
	card, err := ctrl.getCardByIDHandler.Execute(c.Request.Context(), id)
	if err != nil {
		panic(err)
	}

	// Return the response
	c.JSON(http.StatusCreated, gin.H{"data": card})
}
