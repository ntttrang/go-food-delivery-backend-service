package ordermodel

import "time"

// Order represents the orders table.
type Order struct {
	ID         string    `json:"id"`
	UserID     string    `json:"userId"`
	TotalPrice float64   `json:"totalPrice"`
	ShipperID  *string   `json:"shipperId,omitempty"`
	Status     string    `json:"status"`
	CreatedBy  *string   `json:"createdBy,omitempty"`
	UpdatedBy  *string   `json:"updatedBy,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// TableName overrides the table name for Order
func (Order) TableName() string {
	return "orders"
}
