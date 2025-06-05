package ordermodel

import (
	"time"

	"gorm.io/datatypes"
)

// OrderDetail represents the order_details table.
type OrderDetail struct {
	ID         string         `json:"id"`
	OrderID    string         `json:"orderId"`
	FoodOrigin datatypes.JSON `json:"foodOrigin"`
	Price      float64        `json:"price"`
	Quantity   int            `json:"quantity"`
	Discount   float64        `json:"discount"`
	Status     string         `json:"status"`
	CreatedBy  *string        `json:"createdBy,omitempty"`
	UpdatedBy  *string        `json:"updatedBy,omitempty"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
}

// TableName overrides the table name for OrderDetail
func (OrderDetail) TableName() string {
	return "order_details"
}
