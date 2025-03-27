package categorymodel

import (
	"strings"

	"github.com/google/uuid"
)

type CategoryInsertDto struct {
	Name        string `json:"name"`
	Description string `json:"description"`

	Id uuid.UUID `json:"-"`
}

func (c *CategoryInsertDto) Validate() error {
	c.Name = strings.TrimSpace(c.Name)

	if c.Name == "" {
		return ErrNameRequired
	}

	return nil
}

func (c CategoryInsertDto) ConvertToCategory() *Category {
	return &Category{
		Name:        c.Name,
		Description: c.Description,
	}
}
