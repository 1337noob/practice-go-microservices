package repository

import (
	"context"
	"main/services/user/internal/domain"
	"sync"

	"github.com/google/uuid"
)

type InMemoryUserRepository struct {
	users map[string]*domain.User
	mu    sync.Mutex
}

func NewInMemoryRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[string]*domain.User),
	}
}

func (r *InMemoryUserRepository) Create(ctx context.Context, name string, email string, age int) (*domain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user := &domain.User{
		ID:    uuid.NewString(),
		Name:  name,
		Email: email,
		Age:   age,
	}
	r.users[user.ID] = user

	return user, nil
}

func (r *InMemoryUserRepository) Update(ctx context.Context, id string, name string, email string, age int) (*domain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, ok := r.users[id]
	if !ok {
		return nil, domain.ErrUserNotFound
	}

	if name != "" {
		user.Name = name
	}
	if email != "" {
		user.Email = email
	}
	if age != 0 {
		user.Age = age
	}
	r.users[id] = user

	return user, nil
}

func (r *InMemoryUserRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.users[id]
	if !ok {
		return domain.ErrUserNotFound
	}

	delete(r.users, id)

	return nil
}

func (r *InMemoryUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, ok := r.users[id]
	if !ok {
		return nil, domain.ErrUserNotFound
	}

	return user, nil
}

func (r *InMemoryUserRepository) GetAll(ctx context.Context) ([]*domain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var users []*domain.User
	if len(r.users) == 0 {
		return users, nil
	}

	for _, user := range r.users {
		users = append(users, user)
	}

	return users, nil
}
