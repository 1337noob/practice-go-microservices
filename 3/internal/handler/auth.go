package handler

import (
	"context"
	"main/api/proto"
	"main/internal/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthGrpcHandler struct {
	proto.UnimplementedAuthServiceServer
	svs service.AuthServiceInterface
}

func NewAuthGrpcHandler(svs service.AuthServiceInterface) *AuthGrpcHandler {
	return &AuthGrpcHandler{
		svs: svs,
	}
}

func (h *AuthGrpcHandler) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	token, err := h.svs.Login(req.Username, req.Password)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return &proto.LoginResponse{Token: token}, nil
}
