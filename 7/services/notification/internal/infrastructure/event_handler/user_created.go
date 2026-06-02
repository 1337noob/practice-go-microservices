package event_handler

import (
	"errors"
	"fmt"
	"main/pkg/eventbus"
	"main/pkg/eventbus/contracts"
	"main/services/notification/internal/usecase"
)

type UserCreatedHandler struct {
	uc *usecase.SendNotificationUseCase
}

func NewUserCreatedHandler(uc *usecase.SendNotificationUseCase) *UserCreatedHandler {
	return &UserCreatedHandler{uc: uc}
}

func (h *UserCreatedHandler) Handle(event eventbus.Event) error {
	var e, ok = event.(contracts.UserCreated)
	if !ok {
		return errors.New("wrong event type")
	}

	subject := "User Created"
	body := fmt.Sprintf("ID: %s\nName: %s\nEmail: %s\nAge: %d", e.ID, e.Name, e.Email, e.Age)

	return h.uc.Send(e.Email, subject, body)
}
