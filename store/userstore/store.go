package userstore

import (
	"context"

	"github.com/TakumaKurosawa/sqlc-common-transaction/model"
	"github.com/google/uuid"
)

// Store provides user-related database operations
type Store interface {
	CreateUser(ctx context.Context, name, email string) (model.User, error)
	GetUser(ctx context.Context, id uuid.UUID) (model.User, error)
	ListUsers(ctx context.Context) ([]model.User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, name, email string) (model.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}
