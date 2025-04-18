package mediahttpgin

import (
	"context"

	mediaservice "github.com/ntttrang/go-food-delivery-backend-service/modules/media/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type IUploader interface {
	SaveFileUpload(ctx context.Context, fileName string, filePath string, contentType string) error
	GetDomain() string
}

type ICreateCommandHandler interface {
	Execute(ctx context.Context, cmd *mediaservice.CreateCommand) (*uuid.UUID, error)
}

type MediaHTTPController struct {
	createHandler ICreateCommandHandler
	uploader      IUploader
}

func NewMediaHTTPController(
	createHandler ICreateCommandHandler,
	uploader IUploader,
) *MediaHTTPController {
	return &MediaHTTPController{
		createHandler: createHandler,
		uploader:      uploader,
	}
}

func (ctrl *MediaHTTPController) SetupRoutes(router *gin.RouterGroup, auth gin.HandlerFunc) {
	media := router.Group("/medias")
	{
		media.PUT("", auth, ctrl.UploadImageAPI)
	}
}
