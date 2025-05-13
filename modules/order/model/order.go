package ordermodel

import "time"

// Order represents the orders table.
type Order struct {
	ID         string    `gorm:"column:id;primaryKey;type:VARCHAR(36)" json:"id"`
	UserID     string    `gorm:"column:user_id;type:VARCHAR(36);index" json:"userId"`
	TotalPrice float64   `gorm:"column:total_price;type:FLOAT;not null" json:"totalPrice"`
	ShipperID  *string   `gorm:"column:shipper_id;type:VARCHAR(36)"    json:"shipperId,omitempty"`
	Status     string    `gorm:"column:status;type:VARCHAR(50);default:'ACTIVE'" json:"status"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime"      json:"createdAt"`
	UpdatedAt  time.Time `gorm:"column:updated_at;autoUpdateTime"      json:"updatedAt"`
}

// TableName overrides the table name for Order
func (Order) TableName() string {
	return "orders"
}
