package shareinfras

import (
	"gorm.io/gorm"
)

type dbContext struct {
	db *gorm.DB
}

func NewDbContext(db *gorm.DB) IDbContext {
	return &dbContext{
		db: db,
	}
}

func (c *dbContext) GetMainConnection() *gorm.DB {
	return c.db.Session(&gorm.Session{
		NewDB: true,
	})
}
