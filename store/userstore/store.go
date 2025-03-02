//go:generate mockgen -destination=mocks/mock_store.go -package=mocks github.com/TakumaKurosawa/sqlc-common-transaction/store/userstore Store

package userstore

import (
	"context"

	"github.com/TakumaKurosawa/sqlc-common-transaction/model"
	"github.com/google/uuid"
)

// Store defines the interface for user store operations
type Store interface {
	// CreateUser creates a new user
	CreateUser(ctx context.Context, name, email string) (model.User, error)

	// GetUser retrieves a user by ID
	GetUser(ctx context.Context, id uuid.UUID) (model.User, error)

	// ListUsers lists all users
	ListUsers(ctx context.Context) ([]model.User, error)

	// UpdateUser updates a user
	UpdateUser(ctx context.Context, id uuid.UUID, name, email string) (model.User, error)

	// DeleteUser deletes a user
	DeleteUser(ctx context.Context, id uuid.UUID) error
}
