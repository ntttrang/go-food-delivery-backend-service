package categorymodel

import (
	"github.com/google/uuid"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type Category struct {
	Id          uuid.UUID `gorm:"column:id;" json:"id"`
	Name        string    `gorm:"column:name;" json:"name"`
	Description string    `gorm:"column:description;" json:"description"`
	Icon        []byte    `gorm:"column:icon;" json:"icon"`
	Status      string    `gorm:"column:status;" json:"status"`
	sharedModel.DateDto
}

func (Category) TableName() string {
	return "categories"
}
