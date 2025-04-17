package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RPCGetByIdsRequestDTO struct {
	Ids []uuid.UUID `json:"ids"`
}

func (ctl *FoodHttpController) RPCGetByIds(c *gin.Context) {
	var dto RPCGetByIdsRequestDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ids := dto.Ids

	foods, err := ctl.rpcRepo.FindByIds(c.Request.Context(), ids)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": foods})
}
