package sqltransaction

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// txKey is a key for retrieving transaction from context
type txKey struct{}

// Manager implements the transaction.Manager interface using standard SQL
type Manager struct {
	db *sql.DB
}

// New creates a new Manager with the provided database connection
func New(db *sql.DB) *Manager {
	return &Manager{
		db: db,
	}
}

// Begin starts a new transaction
func (m *Manager) Begin(ctx context.Context) (context.Context, error) {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}

	txCtx := context.WithValue(ctx, txKey{}, tx)
	return txCtx, nil
}

// Commit commits the transaction
func (m *Manager) Commit(ctx context.Context) error {
	tx, err := getTx(ctx)
	if err != nil {
		return fmt.Errorf("get transaction: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

// Rollback aborts the transaction
func (m *Manager) Rollback(ctx context.Context) error {
	tx, err := getTx(ctx)
	if err != nil {
		return fmt.Errorf("get transaction: %w", err)
	}

	if err := tx.Rollback(); err != nil {
		return fmt.Errorf("rollback transaction: %w", err)
	}

	return nil
}

// ExecTx executes a function within a transaction
func (m *Manager) ExecTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	txCtx := context.WithValue(ctx, txKey{}, tx)

	if err := fn(txCtx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return fmt.Errorf("transaction failed: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

// getTx extracts the sql.Tx from context
func getTx(ctx context.Context) (*sql.Tx, error) {
	tx, ok := ctx.Value(txKey{}).(*sql.Tx)
	if !ok {
		return nil, errors.New("transaction not found in context")
	}
	return tx, nil
}
