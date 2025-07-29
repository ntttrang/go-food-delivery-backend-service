package usermodel

import (
	"time"

	"github.com/google/uuid"
	sharedmodel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type UserDeviceToken struct {
	Id           uuid.UUID `gorm:"column:id;primaryKey"`
	UserId       uuid.UUID `gorm:"column:user_id;not null;index"`
	Token        string    `gorm:"column:token;not null;unique;size:255"`
	ExpiresAt    time.Time `gorm:"column:expires_at;not null"`
	IsRevoked    bool      `gorm:"column:is_revoked;default:false"`
	IsProduction bool      `gorm:"column:is_production;default:true"`
	OS           string    `gorm:"column:os;size:50"`
	sharedmodel.DateDto
}

func (UserDeviceToken) TableName() string {
	return "user_device_tokens"
}

// IsExpired checks if the device token has expired
func (udt *UserDeviceToken) IsExpired() bool {
	return time.Now().UTC().After(udt.ExpiresAt)
}

// IsValid checks if the device token is valid (not expired and not revoked)
func (udt *UserDeviceToken) IsValid() bool {
	return !udt.IsExpired() && !udt.IsRevoked
}
