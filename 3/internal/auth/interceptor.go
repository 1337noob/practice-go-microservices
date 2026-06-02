package auth

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthInterceptor struct {
	manager       *Manager
	PublicMethods map[string]bool
}

func NewAuthInterceptor(manager *Manager, publicMethods []string) *AuthInterceptor {
	m := make(map[string]bool)
	for _, method := range publicMethods {
		m[method] = true
	}
	return &AuthInterceptor{
		manager:       manager,
		PublicMethods: m,
	}
}

func (i *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if i.PublicMethods[info.FullMethod] {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "metadata is not provided")
		}

		authHeader := md.Get("authorization")
		if len(authHeader) == 0 {
			return nil, status.Error(codes.Unauthenticated, "authorization token is not provided")
		}

		token := strings.TrimPrefix(authHeader[0], "Bearer ")
		if token == authHeader[0] {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		claims, err := i.manager.Validate(token)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		ctx = context.WithValue(ctx, "username", claims.Username)

		return handler(ctx, req)
	}
}
