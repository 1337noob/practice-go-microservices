package main

import (
	"log/slog"
	"main/internal/metrics"
	"net/http"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"log"
	"main/api/proto"
	"main/internal/handler"
	"main/internal/repository"
	"main/internal/service"
	"net"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	logJsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logger := slog.New(logJsonHandler)

	repo := repository.NewInMemoryRepository()
	svs := service.NewUserService(repo)
	srv := handler.NewGrpcHandler(svs)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Error("failed to listen: %v", slog.String("error", err.Error()))
	}

	metricsInterceptor := metrics.NewMetricsInterceptor()
	gs := grpc.NewServer(
		grpc.UnaryInterceptor(metricsInterceptor.Unary(logger)),
	)
	proto.RegisterUserServiceServer(gs, srv)
	reflection.Register(gs)

	http.Handle("/metrics", promhttp.Handler())

	go func() {
		logger.Info("Starting HTTP server on :8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			logger.Error("failed to start HTTP server", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	log.Println("Starting gRPC server on :50051")
	err = gs.Serve(lis)
	if err != nil {
		logger.Error("failed to start gRPC server", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
