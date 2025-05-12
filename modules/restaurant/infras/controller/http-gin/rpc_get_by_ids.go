package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RPCGetByIdsRequestDTO struct {
	Ids []uuid.UUID `json:"ids"`
}

type RPCGetByIdsResponseDTO struct {
	Id               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	Addr             string    `json:"addr"`
	CityId           int       `json:"cityId"`
	Lat              float64   `json:"lat"`
	Lng              float64   `json:"lng"`
	Cover            string    `json:"cover"`
	Logo             string    `json:"logo"`
	ShippingFeePerKm float64   `json:"shippingFeePerKm"`
}

func (ctl *RestaurantHttpController) RPCGetByIds(c *gin.Context) {
	var dto RPCGetByIdsRequestDTO
	var resp []RPCGetByIdsResponseDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ids := dto.Ids

	restaurants, err := ctl.rpcRepo.FindRestaurantByIds(c.Request.Context(), ids)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, r := range restaurants {
		resp = append(resp, RPCGetByIdsResponseDTO{
			Id:               r.Id,
			Name:             r.Name,
			Addr:             r.Addr,
			CityId:           r.CityId,
			Lat:              r.Lat,
			Lng:              r.Lng,
			Cover:            r.Cover,
			Logo:             r.Logo,
			ShippingFeePerKm: r.ShippingFeePerKm,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}
