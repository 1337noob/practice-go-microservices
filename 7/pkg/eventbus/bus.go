package eventbus

import "context"

type EventHandler interface {
	Handle(event Event) error
}

type EventBus interface {
	Subscribe(eventType EventType, handler EventHandler)
	Publish(ctx context.Context, event Event) error
}
