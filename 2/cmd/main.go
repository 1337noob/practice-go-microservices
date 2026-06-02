package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"log"
	"main/api/proto"
	"main/internal/handler"
	"main/internal/repository"
	"main/internal/service"
	"net"
)

func main() {
	repo := repository.NewInMemoryRepository()
	svs := service.NewUserService(repo)
	srv := handler.NewGrpcHandler(svs)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	gs := grpc.NewServer()
	proto.RegisterUserServiceServer(gs, srv)
	reflection.Register(gs)

	log.Println("Starting gRPC server on :50051")
	if err := gs.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
