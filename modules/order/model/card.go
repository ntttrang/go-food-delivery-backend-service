package ordermodel

import "github.com/google/uuid"

type Card struct {
	ID             uuid.UUID `gorm:"column:id" json:"id"`
	Method         string    `gorm:"column:method" json:"method"`
	Provider       string    `gorm:"column:provider" json:"provider"`
	CardholderName string    `gorm:"column:cardholder_name" json:"cardholderName"`
	CardNumber     string    `gorm:"column:card_number" json:"cardNumber"`
	CardType       string    `gorm:"column:card_type" json:"cardType"`
	ExpiryMonth    string    `gorm:"column:expiry_month" json:"expiryMonth"`
	ExpiryYear     string    `gorm:"column:expiry_year" json:"expiryYear"`
	CVV            string    `gorm:"column:cvv" json:"cvv"`
	UserID         uuid.UUID `gorm:"column:user_id" json:"userId"`
	Status         string    `gorm:"column:status" json:"status"`
}
