package shareinfras

import (
	"context"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/ntttrang/go-food-delivery-backend-service/middleware"
	sharecomponent "github.com/ntttrang/go-food-delivery-backend-service/shared/component"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	sharerpc "github.com/ntttrang/go-food-delivery-backend-service/shared/infras/rpc"
)

type IMiddlewareProvider interface {
	Auth() gin.HandlerFunc
}

type IDbContext interface {
	GetMainConnection() *gorm.DB
}

type IUploader interface {
	SaveFileUpload(ctx context.Context, filename string, filePath string, contentType string) error
	GetDomain() string
}

type IAppContext interface {
	MiddlewareProvider() IMiddlewareProvider
	DbContext() IDbContext
	GetConfig() *datatype.Config
	Uploader() IUploader
}

type appContext struct {
	mldProvider IMiddlewareProvider
	dbContext   IDbContext
	config      *datatype.Config
	uploader    IUploader
}

func NewAppContext(db *gorm.DB) IAppContext {
	dbCtx := NewDbContext(db)

	config := datatype.GetConfig()
	introspectRpcClient := sharerpc.NewIntrospectRpcClient(config.UserServiceURL)

	provider := middleware.NewMiddlewareProvider(introspectRpcClient)
	var uploader IUploader
	// Only initialize Minio uploader if the required environment variables are set
	if config.Minio.Domain != "" && config.Minio.AccessKey != "" && config.Minio.SecretKey != "" {
		var err error
		uploader, err = sharecomponent.NewS3Uploader(config.Minio.AccessKey, config.Minio.BucketName, config.Minio.Domain, config.Minio.Region, config.Minio.SecretKey, config.Minio.UseSSL)
		if err != nil {
			panic(err)
		}
	}

	return &appContext{
		mldProvider: provider,
		dbContext:   dbCtx,
		config:      config,
		uploader:    uploader,
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

func (c *appContext) Uploader() IUploader {
	return c.uploader
}
