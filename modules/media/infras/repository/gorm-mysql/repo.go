package gormmysql

import sharedinfras "github.com/ntttrang/go-food-delivery-backend-service/shared/infras"

type ImageRepository struct {
	dbCtx sharedinfras.IDbContext
}

func NewImageRepository(dbCtx sharedinfras.IDbContext) *ImageRepository {
	return &ImageRepository{dbCtx: dbCtx}
}
