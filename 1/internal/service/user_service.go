package service

import (
	"main/internal/dto"
	"main/internal/mapper"
	"main/internal/repository"
)

type Service interface {
	Create(user *dto.CreateUserDto) (*dto.UserDto, error)
	Update(id string, user *dto.UpdateUserDto) (*dto.UserDto, error)
	Delete(id string) error
	GetByID(id string) (*dto.UserDto, error)
	GetAll() ([]*dto.UserDto, error)
}

type UserService struct {
	repo repository.Repository
}

func NewUserService(repo repository.Repository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(user *dto.CreateUserDto) (*dto.UserDto, error) {
	u, err := s.repo.Create(user.Name, user.Email, user.Age)
	if err != nil {
		return nil, err
	}

	return mapper.ModelToDto(u), nil
}

func (s *UserService) Update(id string, user *dto.UpdateUserDto) (*dto.UserDto, error) {
	u, err := s.repo.Update(id, user.Name, user.Email, user.Age)
	if err != nil {
		return nil, err
	}

	return mapper.ModelToDto(u), nil
}

func (s *UserService) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *UserService) GetByID(id string) (*dto.UserDto, error) {
	u, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return mapper.ModelToDto(u), nil
}

func (s *UserService) GetAll() ([]*dto.UserDto, error) {
	users, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	var res []*dto.UserDto
	for _, u := range users {
		res = append(res, mapper.ModelToDto(u))
	}

	return res, nil
}
