package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ntttrang/go-food-delivery-backend-service/modules/food/service"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

func (ctrl *FoodHttpController) CreateFoodCommentAPI(c *gin.Context) {
	var req service.FoodCommentCreateReq

	if err := c.ShouldBindJSON(&req); err != nil {
		panic(datatype.ErrBadRequest.WithError(err.Error()))
	}
	requester := c.MustGet(datatype.KeyRequester).(datatype.Requester)
	req.UserId = requester.Subject()

	if err := ctrl.createCommentFoodCmdHandler.Execute(c.Request.Context(), &req); err != nil {
		panic(err)
	}

	c.JSON(http.StatusCreated, gin.H{"data": req.Id})
}
