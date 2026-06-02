package usecase

type Sender interface {
	Send(to, subject, body string) error
}

type SendNotificationUseCase struct {
	sender Sender
}

func NewSendNotificationUseCase(sender Sender) *SendNotificationUseCase {
	return &SendNotificationUseCase{sender: sender}
}

func (u *SendNotificationUseCase) Send(to, subject, body string) error {
	return u.sender.Send(to, subject, body)
}
