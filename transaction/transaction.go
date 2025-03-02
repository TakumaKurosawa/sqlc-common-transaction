package transaction

import (
	"context"
)

// Manager defines a common interface for transaction management
// This interface allows for implementation across different database drivers
type Manager interface {
	Begin(ctx context.Context) (context.Context, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	ExecTx(ctx context.Context, fn func(ctx context.Context) error) error
}
