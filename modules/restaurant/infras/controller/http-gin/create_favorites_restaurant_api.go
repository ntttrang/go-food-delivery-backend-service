package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	restaurantmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/restaurant/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

// UpdateFavoritesRestaurantAPI: If not exist, add fav restaurant. Otherwise, remove fav restaurant
func (ctrl *RestaurantHttpController) UpdateFavoritesRestaurantAPI(c *gin.Context) {
	var req restaurantmodel.RestaurantLike

	if err := c.ShouldBindJSON(&req); err != nil {
		panic(datatype.ErrBadRequest.WithError(err.Error()))
	}
	requester := c.MustGet(datatype.KeyRequester).(datatype.Requester)
	req.UserID = requester.Subject()

	check, err := ctrl.addFavoritesCmdHdl.Execute(c.Request.Context(), req)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusCreated, gin.H{"data": check})
}
