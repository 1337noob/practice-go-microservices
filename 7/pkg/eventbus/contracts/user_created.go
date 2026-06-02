package contracts

import (
	"main/pkg/eventbus"
	"time"
)

var UserCreatedEventType eventbus.EventType = "user.created"

type UserCreated struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Age       int       `json:"age"`
	Timestamp time.Time `json:"timestamp"`
}

func (e UserCreated) Type() eventbus.EventType {
	return UserCreatedEventType
}
