package handler

import (
	"context"
	"errors"
	"main/api/user/proto"
	"main/services/user/internal/domain"
	"main/services/user/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcHandler struct {
	proto.UnimplementedUserServiceServer
	uc *usecase.UserUseCase
}

func NewGrpcHandler(uc *usecase.UserUseCase) *GrpcHandler {
	return &GrpcHandler{uc: uc}
}

func (h *GrpcHandler) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.User, error) {
	user, err := h.uc.Create(ctx, req.GetName(), req.GetEmail(), int(req.GetAge()))
	if err != nil {
		return nil, err
	}

	return DomainToProto(user), nil
}

func (h *GrpcHandler) UpdateUser(ctx context.Context, req *proto.UpdateUserRequest) (*proto.User, error) {
	user, err := h.uc.Update(ctx, req.GetId(), req.GetName(), req.GetEmail(), int(req.GetAge()))
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, err
	}

	return DomainToProto(user), nil
}

func (h *GrpcHandler) DeleteUser(ctx context.Context, req *proto.DeleteUserRequest) (*proto.DeleteUserResponse, error) {
	err := h.uc.Delete(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, err
	}

	return &proto.DeleteUserResponse{}, nil
}

func (h *GrpcHandler) GetByID(ctx context.Context, req *proto.GetUserByIDRequest) (*proto.User, error) {
	user, err := h.uc.GetByID(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, err
	}

	return DomainToProto(user), nil
}

func (h *GrpcHandler) ListUsers(ctx context.Context, req *proto.ListUsersRequest) (*proto.ListUsersResponse, error) {
	users, err := h.uc.List(ctx)
	if err != nil {
		return nil, err
	}

	var res []*proto.User
	for _, user := range users {
		res = append(res, DomainToProto(user))
	}

	return &proto.ListUsersResponse{
		Users: res,
	}, nil
}
