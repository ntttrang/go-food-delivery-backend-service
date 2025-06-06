package shareinfras

import (
	"context"

	shareComponent "github.com/ntttrang/go-food-delivery-backend-service/shared/component"
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

type IUploader interface {
	SaveFileUpload(ctx context.Context, filename string, filePath string, contentType string) error
	GetDomain() string
}

type IAppContext interface {
	MiddlewareProvider() IMiddlewareProvider
	DbContext() IDbContext
	GetConfig() *datatype.Config
	Uploader() IUploader
	MsgBroker() IMsgBroker
}

type appContext struct {
	mldProvider IMiddlewareProvider
	dbContext   IDbContext
	config      *datatype.Config
	uploader    IUploader
	msgBroker   IMsgBroker
}

func NewAppContext(db *gorm.DB) IAppContext {
	dbCtx := NewDbContext(db)

	config := datatype.GetConfig()
	introspectRpcClient := sharerpc.NewIntrospectRpcClient(config.UserServiceURL)

	provider := middleware.NewMiddlewareProvider(introspectRpcClient)

	var uploader IUploader
	// Only initialize MinioS3 uploader if the required environment variables are set
	if config.MinioS3.Domain != "" && config.MinioS3.AccessKey != "" && config.MinioS3.SecretKey != "" {
		var err error
		uploader, err = shareComponent.NewS3Uploader(config.MinioS3.AccessKey, config.MinioS3.BucketName, config.MinioS3.Domain, config.MinioS3.Region, config.MinioS3.SecretKey, config.MinioS3.UseSSL)
		if err != nil {
			panic(err)
		}
	}

	natsComp := shareComponent.NewNatsComp()
	return &appContext{
		mldProvider: provider,
		dbContext:   dbCtx,
		config:      config,
		uploader:    uploader,
		msgBroker:   natsComp,
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

func (c *appContext) MsgBroker() IMsgBroker {
	return c.msgBroker
}
