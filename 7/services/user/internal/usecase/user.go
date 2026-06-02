package usecase

import (
	"context"
	"main/pkg/eventbus"
	"main/pkg/eventbus/contracts"
	"main/services/user/internal/domain"
	"time"
)

type UserRepository interface {
	Create(ctx context.Context, name string, email string, age int) (*domain.User, error)
	Update(ctx context.Context, id string, name string, email string, age int) (*domain.User, error)
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetAll(ctx context.Context) ([]*domain.User, error)
}

type UserUseCase struct {
	repo UserRepository
	bus  eventbus.EventBus
}

func NewUserUseCase(repo UserRepository, bus eventbus.EventBus) *UserUseCase {
	return &UserUseCase{repo: repo, bus: bus}
}

func (uc *UserUseCase) Create(ctx context.Context, name, email string, age int) (*domain.User, error) {
	user, err := uc.repo.Create(ctx, name, email, age)
	if err != nil {
		return nil, err
	}

	event := contracts.UserCreated{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Age:       user.Age,
		Timestamp: time.Now(),
	}

	err = uc.bus.Publish(ctx, event)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *UserUseCase) Update(ctx context.Context, id, name, email string, age int) (*domain.User, error) {
	user, err := uc.repo.Update(ctx, id, name, email, age)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *UserUseCase) Delete(ctx context.Context, id string) error {
	err := uc.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UserUseCase) GetByID(ctx context.Context, id string) (*domain.User, error) {
	user, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *UserUseCase) List(ctx context.Context) ([]*domain.User, error) {
	users, err := uc.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}
