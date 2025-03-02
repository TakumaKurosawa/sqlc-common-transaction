package transaction

import (
	"context"

	"github.com/TakumaKurosawa/sqlc-common-transaction/store/poststore"
	"github.com/TakumaKurosawa/sqlc-common-transaction/store/userstore"
)

// Manager defines an interface for transaction management across multiple stores
type Manager interface {
	// Begin starts a new transaction
	Begin(ctx context.Context) (context.Context, error)

	// Commit commits the transaction
	Commit(ctx context.Context) error

	// Rollback aborts the transaction
	Rollback(ctx context.Context) error

	// ExecTx executes a function within a transaction
	ExecTx(ctx context.Context, fn func(userStore userstore.Store, postStore poststore.Store) error) error

	// WithTx executes a function using an existing transaction context
	WithTx(ctx context.Context, fn func(userStore userstore.Store, postStore poststore.Store) error) error
}
