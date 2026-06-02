package service

import (
	"errors"
	"main/internal/auth"
)

type AuthServiceInterface interface {
	Login(username, password string) (string, error)
}

type AuthService struct {
	manager *auth.Manager
}

func NewAuthService(manager *auth.Manager) *AuthService {
	return &AuthService{manager: manager}
}

func (s *AuthService) Login(username, password string) (string, error) {
	if username != "admin" || password != "admin" {
		return "", errors.New("invalid credentials")
	}

	token, err := s.manager.Generate(username)
	if err != nil {
		return "", err
	}

	return token, nil
}
