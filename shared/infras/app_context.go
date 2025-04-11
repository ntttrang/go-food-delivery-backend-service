package shareinfras

import (
	sharerpc "github.com/ntttrang/go-food-delivery-backend-service/shared/infras/rpc"

	"github.com/gin-gonic/gin"
	"github.com/ntttrang/go-food-delivery-backend-service/middleware"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	"gorm.io/gorm"
)

type IMiddlewareProvider interface {
	Auth() gin.HandlerFunc
}

type IDbContext interface {
	GetMainConnection() *gorm.DB
}

type IAppContext interface {
	MiddlewareProvider() IMiddlewareProvider
	DbContext() IDbContext
	GetConfig() *datatype.Config
}

type appContext struct {
	mldProvider IMiddlewareProvider
	dbContext   IDbContext
	config      *datatype.Config
}

func NewAppContext(db *gorm.DB) IAppContext {
	dbCtx := NewDbContext(db)

	config := datatype.GetConfig()
	introspectRpcClient := sharerpc.NewIntrospectRpcClient(config.UserServiceURL)

	provider := middleware.NewMiddlewareProvider(introspectRpcClient)

	return &appContext{
		mldProvider: provider,
		dbContext:   dbCtx,
		config:      config,
	}
}

func (c *appContext) MiddlewareProvider() IMiddlewareProvider {
	return c.mldProvider
}

func (c *appContext) DbContext() IDbContext {
	return c.dbContext
}

func (c *appContext) GetConfig() *datatype.Config {
	return datatype.GetConfig()
}
