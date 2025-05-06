package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	foodmodel "github.com/ntttrang/go-food-delivery-backend-service/modules/food/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

func (ctrl *FoodHttpController) SearchFoodAPI(c *gin.Context) {
	// // Check if search functionality is available
	// if ctrl.searchFoodQueryHandler == nil {
	// 	c.JSON(http.StatusServiceUnavailable, gin.H{
	// 		"error": "Search functionality is not available. Elasticsearch is not configured.",
	// 	})
	// 	return
	// }

	// var searchQuery foodmodel.FoodSearchQuery
	// var pagingDto sharedModel.PagingDto
	// var sortingDto sharedModel.SortingDto

	// // Bind query parameters
	// if err := c.ShouldBind(&searchQuery); err != nil {
	// 	panic(datatype.ErrBadRequest.WithError(err.Error()))
	// }

	// if err := c.ShouldBind(&pagingDto); err != nil {
	// 	panic(datatype.ErrBadRequest.WithError(err.Error()))
	// }

	// if err := c.ShouldBind(&sortingDto); err != nil {
	// 	panic(datatype.ErrBadRequest.WithError(err.Error()))
	// }

	// // Process pagination
	// pagingDto.Process()

	// // Create search request
	// req := foodmodel.FoodSearchReq{
	// 	FoodSearchQuery: searchQuery,
	// }
	// req.PagingDto = pagingDto

	// Create a request object to bind the JSON body
	var req foodmodel.FoodSearchReq

	// Bind JSON body
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(datatype.ErrBadRequest.WithError(err.Error()))
	}

	// Process pagination
	req.Process()

	// Execute search
	result, err := ctrl.searchFoodQueryHandler.Execute(c.Request.Context(), req)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (ctrl *FoodHttpController) ReindexFoodAPI(c *gin.Context) {
	// Check if search functionality is available
	if ctrl.syncFoodIndexCommandHandler == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Reindex functionality is not available. Elasticsearch is not configured.",
		})
		return
	}

	// This endpoint should be protected and only accessible by admins
	err := ctrl.syncFoodIndexCommandHandler.ReindexAll(c.Request.Context())
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Food index reindexed successfully"})
}
