package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	foodmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/food/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

func (ctrl *FoodHttpController) ListFavoriteFoodAPI(c *gin.Context) {
	var pagingDto sharedModel.PagingDto
	if err := c.ShouldBind(&pagingDto); err != nil {
		panic(datatype.ErrBadRequest.WithError(err.Error()))
	}

	pagingDto.Process()

	requester := c.MustGet(datatype.KeyRequester).(datatype.Requester)
	req := foodmodel.FavoriteFoodListReq{
		UserId:    requester.Subject(),
		PagingDto: pagingDto,
	}

	result, err := ctrl.favoriteFoodQueryHdl.Execute(c.Request.Context(), req)

	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}
