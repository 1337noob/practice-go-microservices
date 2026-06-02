package main

import (
	"main/internal/auth"
	"time"

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
	authManager := auth.NewManager("secret", time.Minute)
	authSvc := service.NewAuthService(authManager)
	authSrv := handler.NewAuthGrpcHandler(authSvc)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	authInterceptor := auth.NewAuthInterceptor(authManager, []string{
		"/proto.AuthService/Login",
	})

	gs := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor.Unary()),
	)

	proto.RegisterUserServiceServer(gs, srv)
	proto.RegisterAuthServiceServer(gs, authSrv)
	reflection.Register(gs)

	log.Println("Starting gRPC server on :50051")
	if err := gs.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
