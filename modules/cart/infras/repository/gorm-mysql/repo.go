package cartgormmysql

import shareinfras "github.com/ntttrang/go-food-delivery-backend-service/shared/infras"

type CartRepo struct {
	dbCtx shareinfras.IDbContext
}

func NewCartRepo(dbCtx shareinfras.IDbContext) *CartRepo {
	return &CartRepo{dbCtx: dbCtx}
}
