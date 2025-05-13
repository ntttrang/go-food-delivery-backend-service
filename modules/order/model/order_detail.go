package ordermodel

import (
	"time"

	"gorm.io/datatypes"
)

// OrderDetail represents the order_details table.
type OrderDetail struct {
	ID         string         `gorm:"column:id;primaryKey;type:VARCHAR(36)" json:"id"`
	OrderID    string         `gorm:"column:order_id;type:VARCHAR(36);index"    json:"orderId"`
	FoodOrigin datatypes.JSON `gorm:"column:food_origin;type:JSON"             json:"foodOrigin"`
	Price      float64        `gorm:"column:price;type:FLOAT;not null"         json:"price"`
	Quantity   int            `gorm:"column:quantity;type:INT;not null"        json:"quantity"`
	Discount   float64        `gorm:"column:discount;type:FLOAT;default:0"      json:"discount"`
	Status     string         `gorm:"column:status;type:VARCHAR(50);default:'ACTIVE'" json:"status"`
	CreatedAt  time.Time      `gorm:"column:created_at;autoCreateTime"         json:"createdAt"`
	UpdatedAt  time.Time      `gorm:"column:updated_at;autoUpdateTime"         json:"updatedAt"`
}

// TableName overrides the table name for OrderDetail
func (OrderDetail) TableName() string {
	return "order_details"
}
