package main

import (
	"context"
	"log"
	"main/pkg/eventbus/contracts"
	"main/pkg/eventbus/kafka"
	"main/services/notification/internal/infrastructure/event_handler"
	"main/services/notification/internal/infrastructure/smtp"
	"main/services/notification/internal/usecase"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	brokers := []string{"localhost:9092"}
	consumerGroup := "notification-group"
	sender := smtp.NewSmtpSender("localhost", "1025", "noreply@example.com")
	bus, err := kafka.NewKafkaEventBus(brokers)
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

	err = bus.StartConsumerGroup(ctxWithCancel, brokers, consumerGroup)
	if err != nil {
		log.Println("Consumer group error: ", err)
	}
}
