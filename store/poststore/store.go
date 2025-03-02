package poststore

import (
	"context"

	"github.com/TakumaKurosawa/sqlc-common-transaction/model"
	"github.com/google/uuid"
)

// Store provides post-related database operations
type Store interface {
	CreatePost(ctx context.Context, userID uuid.UUID, title, content string) (model.Post, error)
	GetPost(ctx context.Context, id uuid.UUID) (model.Post, error)
	ListPostsByUser(ctx context.Context, userID uuid.UUID) ([]model.Post, error)
	UpdatePost(ctx context.Context, id uuid.UUID, title, content string) (model.Post, error)
	DeletePost(ctx context.Context, id uuid.UUID) error
}
