package foodmodel

import (
	"strings"

	"github.com/google/uuid"
)

type FoodInsertDto struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`

	Id uuid.UUID `json:"-"`
}

func (c *FoodInsertDto) Validate() error {
	c.Name = strings.TrimSpace(c.Name)

	if c.Name == "" {
		return ErrNameRequired
	}

	return nil
}

func (c FoodInsertDto) ConvertToFood() *Food {
	return &Food{
		Name:        c.Name,
		Description: c.Description,
		Price:       c.Price,
	}
}
