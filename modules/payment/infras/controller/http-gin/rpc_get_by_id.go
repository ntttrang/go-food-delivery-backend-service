package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RPCGetByIdsRequestDTO struct {
	Id uuid.UUID `json:"id"`
}

func (ctrl *CardController) RPCGetById(c *gin.Context) {
	var dto RPCGetByIdsRequestDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := dto.Id

	card, err := ctrl.repo.FindByID(c.Request.Context(), id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": card})
}
