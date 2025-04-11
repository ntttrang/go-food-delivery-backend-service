package repository

import shareinfras "github.com/ntttrang/go-food-delivery-backend-service/shared/infras"

type CategoryRepo struct {
	dbCtx shareinfras.IDbContext
}

func NewCategoryRepo(dbCtx shareinfras.IDbContext) *CategoryRepo {
	return &CategoryRepo{dbCtx: dbCtx}
}
