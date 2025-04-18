package categorymodel

import (
	"github.com/google/uuid"
	sharedModel "github.com/ntttrang/go-food-delivery-backend-service/shared/model"
)

type CategoryUpdateReq struct {
	// Use pointer to accept empty string
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Status      *string `json:"status"`

	Id uuid.UUID `json:"-"`
}

func (CategoryUpdateReq) TableName() string {
	return Category{}.TableName()
}

func (c CategoryUpdateReq) Validate() error {
	if c.Status != nil && *c.Status != sharedModel.StatusActive && *c.Status != sharedModel.StatusDelete && *c.Status != sharedModel.StatusInactive {
		return ErrCategoryStatusInvalid
	}
	return nil
}
