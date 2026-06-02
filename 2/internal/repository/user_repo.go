package repository

import (
	"main/internal/errors"
	"main/internal/model"
	"sync"

	"github.com/google/uuid"
)

type Repository interface {
	Create(name string, email string, age int) (*model.User, error)
	Update(id string, name string, email string, age int) (*model.User, error)
	Delete(id string) error
	GetByID(id string) (*model.User, error)
	GetAll() ([]*model.User, error)
}

type InMemoryUserRepository struct {
	users map[string]*model.User
	mu    sync.Mutex
}

func NewInMemoryRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[string]*model.User),
	}
}

func (r *InMemoryUserRepository) Create(name string, email string, age int) (*model.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user := &model.User{
		ID:    uuid.NewString(),
		Name:  name,
		Email: email,
		Age:   age,
	}
	r.users[user.ID] = user

	return user, nil
}

func (r *InMemoryUserRepository) Update(id string, name string, email string, age int) (*model.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, ok := r.users[id]
	if !ok {
		return nil, errors.ErrUserNotFound
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

func (r *InMemoryUserRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.users[id]
	if !ok {
		return errors.ErrUserNotFound
	}

	delete(r.users, id)

	return nil
}

func (r *InMemoryUserRepository) GetByID(id string) (*model.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, ok := r.users[id]
	if !ok {
		return nil, errors.ErrUserNotFound
	}

	return user, nil
}

func (r *InMemoryUserRepository) GetAll() ([]*model.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var users []*model.User
	if len(r.users) == 0 {
		return users, nil
	}

	for _, user := range r.users {
		users = append(users, user)
	}

	return users, nil
}
