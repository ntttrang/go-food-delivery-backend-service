package shareinfras

import (
	"context"

	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
)

type IMsgBroker interface {
	Publish(ctx context.Context, topic string, evt *datatype.AppEvent) error
}
