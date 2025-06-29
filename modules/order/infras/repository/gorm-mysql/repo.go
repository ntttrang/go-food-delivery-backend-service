package ordergormmysql

import shareinfras "github.com/ntttrang/go-food-delivery-backend-service/shared/infras"

type OrderRepo struct {
	dbCtx shareinfras.IDbContext
}

func NewOrderRepo(dbCtx shareinfras.IDbContext) *OrderRepo {
	return &OrderRepo{dbCtx: dbCtx}
}
