package handler

import (
	"context"
	"errors"
	"main/api/proto"
	"main/internal/dto"
	errors2 "main/internal/errors"
	"main/internal/mapper"
	"main/internal/service"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcHandler struct {
	proto.UnimplementedUserServiceServer
	svs service.Service
}

func NewGrpcHandler(svs service.Service) *GrpcHandler {
	return &GrpcHandler{svs: svs}
}

func (h *GrpcHandler) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.User, error) {
	createUser := &dto.CreateUserDto{
		Name:  req.GetName(),
		Email: req.GetEmail(),
		Age:   int(req.GetAge()),
	}

	user, err := h.svs.Create(createUser)
	if err != nil {
		return nil, err
	}

	return mapper.DtoToProto(user), nil
}

func (h *GrpcHandler) UpdateUser(ctx context.Context, req *proto.UpdateUserRequest) (*proto.User, error) {
	updateUser := &dto.UpdateUserDto{
		Name:  req.GetName(),
		Email: req.GetEmail(),
		Age:   int(req.GetAge()),
	}

	user, err := h.svs.Update(req.Id, updateUser)
	if err != nil {
		if errors.Is(err, errors2.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, err
	}

	return mapper.DtoToProto(user), nil
}

func (h *GrpcHandler) DeleteUser(ctx context.Context, req *proto.DeleteUserRequest) (*proto.DeleteUserResponse, error) {
	err := h.svs.Delete(req.Id)
	if err != nil {
		if errors.Is(err, errors2.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, err
	}

	return &proto.DeleteUserResponse{}, nil
}

func (h *GrpcHandler) GetByID(ctx context.Context, req *proto.GetUserByIDRequest) (*proto.User, error) {
	user, err := h.svs.GetByID(req.Id)
	if err != nil {
		if errors.Is(err, errors2.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, err
	}

	return mapper.DtoToProto(user), nil
}

func (h *GrpcHandler) ListUsers(ctx context.Context, req *proto.ListUsersRequest) (*proto.ListUsersResponse, error) {
	users, err := h.svs.GetAll()
	if err != nil {
		return nil, err
	}

	var res []*proto.User
	for _, user := range users {
		res = append(res, mapper.DtoToProto(user))
	}

	// test metrics
	time.Sleep(100 * time.Millisecond)

	return &proto.ListUsersResponse{
		Users: res,
	}, nil
}
