package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	service "github.com/ntttrang/go-food-delivery-backend-service/modules/user/service"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

func (ctrl *UserHttpController) GetUserDetailAPI(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		panic(datatype.ErrBadRequest.WithError(err.Error()))
	}

	restaurant, err := ctrl.getDetailQueryHdl.Execute(c.Request.Context(), service.UserDetailReq{Id: id})

	if err != nil {
		// c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// return
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"data": restaurant})
}
