package main

import (
	"log"
	"main/api/user/proto"
	"main/pkg/eventbus/kafka"
	handler "main/services/user/internal/infrastructure/grpc"
	"main/services/user/internal/infrastructure/repository"
	"main/services/user/internal/usecase"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	brokers := []string{"localhost:9092"}

	repo := repository.NewInMemoryRepository()
	bus, err := kafka.NewKafkaEventBus(brokers)
	if err != nil {
		log.Fatal(err)
	}
	uc := usecase.NewUserUseCase(repo, bus)
	srv := handler.NewGrpcHandler(uc)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	gs := grpc.NewServer()
	proto.RegisterUserServiceServer(gs, srv)
	reflection.Register(gs)

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-signals
		log.Println("Shutting down gracefully")
		gs.GracefulStop()
		log.Println("gRPC server stopped")
	}()

	log.Println("Starting gRPC server on :50051")
	if err := gs.Serve(lis); err != nil {
		log.Fatal(err)
	}

	log.Println("Shutting down event bus")
	err = bus.Close()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("User service stopped gracefully")
}
