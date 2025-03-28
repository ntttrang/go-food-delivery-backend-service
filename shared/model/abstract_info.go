package sharedModel

import (
	"time"
)

type AbstractInfo struct {
	CreatedAt *time.Time `json:"createdAt" gorm:"column:created_at;"`
	UpdatedAt *time.Time `json:"updatedAt" gorm:"column:updated_at;"`
}
