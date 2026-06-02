package main

import (
	"fmt"
	"log"
	"main/api/user/proto"
	"main/pkg/eventbus/kafka"
	"main/services/user/internal/config"
	handler "main/services/user/internal/infrastructure/grpc"
	"main/services/user/internal/infrastructure/repository"
	"main/services/user/internal/usecase"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.Load()

	repo := repository.NewInMemoryRepository()
	bus, err := kafka.NewKafkaEventBus([]string{cfg.KafkaBroker})
	if err != nil {
		log.Fatal(err)
	}
	uc := usecase.NewUserUseCase(repo, bus)
	srv := handler.NewGrpcHandler(uc)

	addr := fmt.Sprintf(":%d", cfg.GRPCPort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	gs := grpc.NewServer()
	proto.RegisterUserServiceServer(gs, srv)
	reflection.Register(gs)

	log.Printf("Starting gRPC server on %s", addr)
	if err := gs.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
