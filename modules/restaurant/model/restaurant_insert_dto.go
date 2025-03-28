package restaurantmodel

import (
	"encoding/json"
	"strings"

	"github.com/google/uuid"
)

type RestaurantInsertDto struct {
	OwnerId string  `json:"ownerId"`
	Name    string  `json:"name"`
	Addr    string  `json:"addr"`
	CityId  int     `json:"cityId"`
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`
	//Cover            string  `json:"cover"` // json
	//Logo             string  `json:"logo"`  // json
	ShippingFeePerKm float64             `json:"shippingFeePerKm"`
	Foods            []RestaurantFoodDto `json:"foods"`

	Id uuid.UUID `json:"-"` // Internal BE
}

func (r RestaurantInsertDto) Validate() error {
	r.Name = strings.TrimSpace(r.Name)

	if r.Name == "" {
		return ErrNameRequired
	}

	return nil
}

func (r RestaurantInsertDto) ConvertToRestaurant() *Restaurant {
	return &Restaurant{
		OwnerId:          r.OwnerId,
		Name:             r.Name,
		Addr:             r.Addr,
		CityId:           r.CityId,
		Lat:              r.Lat,
		Lng:              r.Lng,
		ShippingFeePerKm: r.ShippingFeePerKm,

		// TODO: Hard code
		Cover: json.RawMessage(`{"key": "value"}`),
		Logo:  json.RawMessage(`{"key": "value"}`),
	}
}
