package mediaservice

import (
	"context"
	"time"

	"github.com/google/uuid"

	mediamodel "github.com/ntttrang/go-food-delivery-backend-service/modules/media/model"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

type CreateCommand struct {
	ImageCreate mediamodel.ImageCreateDTO
}

type ICreateRepo interface {
	Insert(ctx context.Context, data *mediamodel.Image) error
}

type CreateCommandHandler struct {
	mediaRepo ICreateRepo
}

func NewCreateCommandHandler(mediaRepo ICreateRepo) *CreateCommandHandler {
	return &CreateCommandHandler{mediaRepo: mediaRepo}
}

func (hdl *CreateCommandHandler) Execute(ctx context.Context, cmd *CreateCommand) (*uuid.UUID, error) {
	newId, _ := uuid.NewV7()
	now := time.Now().UTC()

	media := mediamodel.Image{
		Id:        newId,
		Filename:  cmd.ImageCreate.Filename,
		CloudName: cmd.ImageCreate.CloudName,
		Size:      cmd.ImageCreate.Size,
		Ext:       cmd.ImageCreate.Ext,
		Status:    string(datatype.StatusActive),
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := hdl.mediaRepo.Insert(ctx, &media); err != nil {
		return nil, datatype.ErrInternalServerError.WithWrap(err).WithDebug(err.Error())
	}

	return &media.Id, nil
}
