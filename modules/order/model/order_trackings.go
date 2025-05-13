package ordermodel

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// OrderTracking represents the order_trackings table.
type OrderTracking struct {
	ID              string         `gorm:"column:id;primaryKey;type:VARCHAR(36)" json:"id"`
	OrderID         string         `gorm:"column:order_id;type:VARCHAR(36);index" json:"orderId"`
	State           string         `gorm:"column:state;type:ENUM('waiting_for_shipper','preparing','on_the_way','delivered','cancel');not null" json:"state"`
	PaymentStatus   string         `gorm:"column:payment_status;type:ENUM('pending','paid');not null" json:"paymentStatus"`
	PaymentMethod   string         `gorm:"column:payment_method;type:ENUM('cash','card');not null" json:"paymentMethod"`
	CardId          uuid.UUID      `gorm:"column:card_id;" json:"cardId"`
	DeliveryAddress datatypes.JSON `gorm:"column:delivery_address;type:JSON"           json:"deliveryAddress"`
	DeliveryFee     float64        `gorm:"column:delivery_fee;type:FLOAT;default:0"     json:"deliveryFee"`
	EstimatedTime   int            `gorm:"column:estimated_time;type:INT;default:0"     json:"estimatedTime"`
	DeliveryTime    int            `gorm:"column:delivery_time;type:INT;default:0"      json:"deliveryTime"`
	RestaurantID    string         `gorm:"column:restaurant_id;type:VARCHAR(36);not null" json:"restaurantId"`
	Status          string         `gorm:"column:status;type:VARCHAR(50);default:'ACTIVE'" json:"status"`
	CreatedAt       time.Time      `gorm:"column:created_at;autoCreateTime"             json:"createdAt"`
	UpdatedAt       time.Time      `gorm:"column:updated_at;autoUpdateTime"             json:"updatedAt"`
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
