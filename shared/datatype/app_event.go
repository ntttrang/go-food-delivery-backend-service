package datatype

type AppEvent struct {
	Topic string
	Data  interface{}
}

type AppEventOpt func(*AppEvent)

func WithTopic(topic string) AppEventOpt {
	return func(evt *AppEvent) {
		evt.Topic = topic
	}
}

func WithData(data interface{}) AppEventOpt {
	return func(evt *AppEvent) {
		evt.Data = data
	}
}

func NewAppEvent(opts ...AppEventOpt) *AppEvent {
	evt := &AppEvent{}

	for _, opt := range opts {
		opt(evt)
	}

	return evt
}
