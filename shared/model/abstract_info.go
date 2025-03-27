package sharedModel

import (
	"time"

	"gorm.io/gorm"
)

type AbstractInfo struct {
	CreatedAt *time.Time `json:"createdAt" gorm:"column:created_at;"`
	UpdatedAt *time.Time `json:"updatedAt" gorm:"column:updated_at;"`
}

func (c *AbstractInfo) BeforeCreate(tx *gorm.DB) {
	currentTime := time.Now().UTC()
	c.CreatedAt = &currentTime
	c.UpdatedAt = &currentTime
}
