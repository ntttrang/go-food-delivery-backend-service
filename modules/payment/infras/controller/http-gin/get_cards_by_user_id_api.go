package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

// GetCardsByUserID handles the get cards by user ID request
func (ctrl *CardController) GetCardsByUserID(c *gin.Context) {
	// Parse the user ID from the URL
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		panic(datatype.ErrBadRequest.WithError(err.Error()))
	}

	// Execute the query
	cards, err := ctrl.getCardsByUserIDHandler.Execute(c.Request.Context(), userID)
	if err != nil {
		panic(err)
	}

	// Return the response
	c.JSON(http.StatusCreated, gin.H{"data": gin.H{"items": cards}})
}
