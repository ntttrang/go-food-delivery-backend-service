package gormmysql

import shareinfras "github.com/ntttrang/go-food-delivery-backend-service/shared/infras"

type CardRepo struct {
	dbCtx shareinfras.IDbContext
}

func NewCardRepo(dbCtx shareinfras.IDbContext) *CardRepo {
	return &CardRepo{dbCtx: dbCtx}
}
