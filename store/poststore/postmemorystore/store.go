package postmemorystore

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/TakumaKurosawa/sqlc-common-transaction/model"
	"github.com/TakumaKurosawa/sqlc-common-transaction/store/poststore"
	"github.com/google/uuid"
)

var (
	ErrPostNotFound = errors.New("post not found")
)

type memoryStore struct {
	mu    sync.RWMutex
	posts map[uuid.UUID]model.Post
}

// New creates a new in-memory implementation of poststore.Store
func New() poststore.Store {
	return &memoryStore{
		posts: make(map[uuid.UUID]model.Post),
	}
}

func (s *memoryStore) CreatePost(_ context.Context, userID uuid.UUID, title, content string) (model.Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	post := model.Post{
		ID:        uuid.New(),
		UserID:    userID,
		Title:     title,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
	}

	s.posts[post.ID] = post
	return post, nil
}

func (s *memoryStore) GetPost(_ context.Context, id uuid.UUID) (model.Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	post, exists := s.posts[id]
	if !exists {
		return model.Post{}, ErrPostNotFound
	}

	return post, nil
}

func (s *memoryStore) ListPostsByUser(_ context.Context, userID uuid.UUID) ([]model.Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []model.Post
	for _, post := range s.posts {
		if post.UserID == userID {
			result = append(result, post)
		}
	}

	return result, nil
}

func (s *memoryStore) UpdatePost(_ context.Context, id uuid.UUID, title, content string) (model.Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	post, exists := s.posts[id]
	if !exists {
		return model.Post{}, ErrPostNotFound
	}

	post.Title = title
	post.Content = content
	post.UpdatedAt = time.Now()

	s.posts[id] = post
	return post, nil
}

func (s *memoryStore) DeletePost(_ context.Context, id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.posts[id]; !exists {
		return ErrPostNotFound
	}

	delete(s.posts, id)
	return nil
}
