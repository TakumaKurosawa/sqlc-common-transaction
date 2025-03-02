package usermemorystore

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/TakumaKurosawa/sqlc-common-transaction/model"
	"github.com/TakumaKurosawa/sqlc-common-transaction/store/userstore"
	"github.com/google/uuid"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type memoryStore struct {
	mu    sync.RWMutex
	users map[uuid.UUID]model.User
}

// New creates a new in-memory implementation of userstore.Store
func New() userstore.Store {
	return &memoryStore{
		users: make(map[uuid.UUID]model.User),
	}
}

func (s *memoryStore) CreateUser(_ context.Context, name, email string) (model.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	user := model.User{
		ID:        uuid.New(),
		Name:      name,
		Email:     email,
		CreatedAt: now,
		UpdatedAt: now,
	}

	s.users[user.ID] = user
	return user, nil
}

func (s *memoryStore) GetUser(_ context.Context, id uuid.UUID) (model.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[id]
	if !exists {
		return model.User{}, ErrUserNotFound
	}

	return user, nil
}

func (s *memoryStore) ListUsers(_ context.Context) ([]model.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make([]model.User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}

	return users, nil
}

func (s *memoryStore) UpdateUser(_ context.Context, id uuid.UUID, name, email string) (model.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, exists := s.users[id]
	if !exists {
		return model.User{}, ErrUserNotFound
	}

	user.Name = name
	user.Email = email
	user.UpdatedAt = time.Now()

	s.users[id] = user
	return user, nil
}

func (s *memoryStore) DeleteUser(_ context.Context, id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[id]; !exists {
		return ErrUserNotFound
	}

	delete(s.users, id)
	return nil
}
