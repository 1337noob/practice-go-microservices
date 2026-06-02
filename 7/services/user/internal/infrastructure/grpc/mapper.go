package handler

import (
	"main/api/user/proto"
	"main/services/user/internal/domain"
)

func DomainToProto(user *domain.User) *proto.User {
	return &proto.User{
		Id:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Age:   int32(user.Age),
	}
}
