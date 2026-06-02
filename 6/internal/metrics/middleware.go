package metrics

import (
	"context"
	"fmt"
	"time"

	"log/slog"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

var (
	RequestTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "grpc_requests_total",
		Help: "The total number of grpc requests",
	}, []string{"method", "code"})

	RequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "grpc_request_duration_seconds",
		Help: "The grpc request duration in seconds",
	}, []string{"method", "code"})

	RequestsInFlight = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "grpc_requests_in_flight",
		Help: "The current number of grpc requests in flight",
	})
)

type MetricsInterceptor struct {
	PublicMethods map[string]bool
}

func NewMetricsInterceptor() *MetricsInterceptor {
	return &MetricsInterceptor{}
}

func (i *MetricsInterceptor) Unary(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		RequestsInFlight.Inc()
		start := time.Now()

		res, err := handler(ctx, req)

		duration := time.Since(start)
		RequestsInFlight.Dec()

		method := info.FullMethod
		code := status.Code(err).String()

		RequestTotal.WithLabelValues(method, code).Inc()
		RequestDuration.WithLabelValues(method, code).Observe(duration.Seconds())

		logger.InfoContext(ctx, "GRPC Request: ",
			slog.String("method", method),
			slog.String("code", code),
			slog.String("duration", fmt.Sprintf("%v", duration.String())),
		)

		return res, err
	}
}
