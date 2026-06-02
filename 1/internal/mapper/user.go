package mapper

import (
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
