//go:generate mockgen -destination=mocks/mock_store.go -package=mocks github.com/TakumaKurosawa/sqlc-common-transaction/store/poststore Store

package poststore

import (
	"context"

	"github.com/TakumaKurosawa/sqlc-common-transaction/model"
	"github.com/google/uuid"
)

// Store defines the interface for post store operations
type Store interface {
	// CreatePost creates a new post
	CreatePost(ctx context.Context, userID uuid.UUID, title, content string) (model.Post, error)

	// GetPost retrieves a post by ID
	GetPost(ctx context.Context, id uuid.UUID) (model.Post, error)

	// ListPostsByUser lists all posts by a user
	ListPostsByUser(ctx context.Context, userID uuid.UUID) ([]model.Post, error)

	// UpdatePost updates a post
	UpdatePost(ctx context.Context, id uuid.UUID, title, content string) (model.Post, error)

	// DeletePost deletes a post
	DeletePost(ctx context.Context, id uuid.UUID) error
}
