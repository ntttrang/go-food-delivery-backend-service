package restaurantmodel

type RestaurantUpdateDto struct {
	Name   *string  `json:"name"`
	Addr   *string  `json:"addr"`
	CityId *int     `json:"cityId"`
	Lat    *float64 `json:"lat"`
	Lng    *float64 `json:"lng"`
	// Cover            *json.RawMessage `json:"cover"`
	// Logo             *json.RawMessage `json:"logo"`
	ShippingFeePerKm *float64 `json:"shippingFeePerKm"`
	Status           *string  `json:"status"`
}
