package restaurantmodel

type RestaurantSearchDto struct {
	// CategoryName *string `json:"categoryName"`
	OwnerId *string `json:"ownerId" form:"ownerId"`
	CityId  *int    `json:"cityId" form:"cityId"`
	Status  *string `json:"status" form:"status"`
}
