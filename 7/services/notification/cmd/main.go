package main

import (
	"context"
	"log"
	"main/pkg/eventbus/contracts"
	"main/pkg/eventbus/kafka"
	"main/services/notification/internal/config"
	"main/services/notification/internal/infrastructure/event_handler"
	"main/services/notification/internal/infrastructure/smtp"
	"main/services/notification/internal/usecase"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.Load()

	sender := smtp.NewSmtpSender(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPFrom)
	bus, err := kafka.NewKafkaEventBus([]string{cfg.KafkaBroker})
	if err != nil {
		log.Fatal(err)
	}
	uc := usecase.NewSendNotificationUseCase(sender)
	onUserCreated := event_handler.NewUserCreatedHandler(uc)
	bus.Subscribe(contracts.UserCreatedEventType, onUserCreated)

	ctx := context.Background()
	ctxWithCancel, cancel := context.WithCancel(ctx)

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-signals
		cancel()
		log.Println("Shutting down gracefully")
	}()

	err = bus.StartConsumerGroup(ctxWithCancel, []string{cfg.KafkaBroker}, cfg.ConsumerGroup)
	if err != nil {
		log.Println("Consumer group error: ", err)
	}
}
