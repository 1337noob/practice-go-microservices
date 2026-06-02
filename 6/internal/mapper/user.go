package mapper

import (
	"main/api/proto"
	"main/internal/dto"
	"main/internal/model"
)

func ModelToDto(user *model.User) *dto.UserDto {
	return &dto.UserDto{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Age:   user.Age,
	}
}

func DtoToProto(user *dto.UserDto) *proto.User {
	return &proto.User{
		Id:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Age:   int32(user.Age),
	}
}
