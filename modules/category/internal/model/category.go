package categorymodel

import (
	"github.com/google/uuid"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type Category struct {
	Id          uuid.UUID `gorm:"column:id;"`
	Name        string    `gorm:"column:name;"`
	Description string    `gorm:"column:description;"`
	Icon        []byte    `gorm:"column:icon;"`
	Status      string    `gorm:"column:status;"`
	sharedModel.AbstractInfo
}
