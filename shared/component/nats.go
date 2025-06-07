package shareComponent

import (
	"context"
	"encoding/json"
	"log"

	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	"go.opentelemetry.io/otel"

	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
)

type natsComp struct {
	nc *nats.Conn
}

func NewNatsComp() *natsComp {
	nc, err := nats.Connect(datatype.GetConfig().NatsURL)
	if err != nil {
		log.Fatal(err)
	}

	return &natsComp{nc: nc}
}

func (c *natsComp) Publish(ctx context.Context, topic string, evt *datatype.AppEvent) error {
	_, dbSpanPlbNoti := otel.Tracer("").Start(ctx, "publish-msg")
	defer dbSpanPlbNoti.End()

	dataByte, err := json.Marshal(evt.Data)

	if err != nil {
		return errors.WithStack(err)
	}

	return c.nc.Publish(topic, dataByte)
}
