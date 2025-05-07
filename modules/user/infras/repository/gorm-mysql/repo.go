package usergormmysql

import (
	shareinfras "github.com/ntttrang/go-food-delivery-backend-service/shared/infras"
)

type UserRepo struct {
	dbCtx shareinfras.IDbContext
}

func NewUserRepo(dbCtx shareinfras.IDbContext) *UserRepo {
	return &UserRepo{dbCtx: dbCtx}
}

type UserAddressRepo struct {
	dbCtx shareinfras.IDbContext
}

func NewUserAddressRepo(dbCtx shareinfras.IDbContext) *UserAddressRepo {
	return &UserAddressRepo{dbCtx: dbCtx}
}
