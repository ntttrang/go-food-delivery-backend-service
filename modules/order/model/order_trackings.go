package ordermodel

import (
	"time"

	"gorm.io/datatypes"
)

// OrderTracking represents the order_trackings table.
type OrderTracking struct {
	ID              string         `json:"id"`
	OrderID         string         `json:"orderId"`
	State           string         `json:"state"`
	CancelReason    string         `json:"cancelReason"`
	PaymentStatus   string         `json:"paymentStatus"`
	PaymentMethod   string         `json:"paymentMethod"`
	CardId          *string        `json:"cardId"`
	DeliveryAddress datatypes.JSON `json:"deliveryAddress"`
	DeliveryFee     float64        `json:"deliveryFee"`
	EstimatedTime   int            `json:"estimatedTime"`
	DeliveryTime    int            `json:"deliveryTime"`
	RestaurantID    string         `json:"restaurantId"`
	Status          string         `json:"status"`
	CreatedBy       *string        `json:"createdBy,omitempty"`
	UpdatedBy       *string        `json:"updatedBy,omitempty"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
}

// TableName overrides the table name for OrderTracking
func (OrderTracking) TableName() string {
	return "order_trackings"
}

type Address struct {
	CityId   int      `json:"cityId"`
	CityName string   `json:"cityName"`
	Addr     string   `json:"addr"`
	Lat      *float64 `json:"lat"`
	Lng      *float64 `json:"lng"`
}
