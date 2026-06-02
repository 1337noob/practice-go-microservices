package eventbus

type EventType string

type Event interface {
	Type() EventType
}
