package pgxtransaction

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// pgxTxKey is a key for retrieving transaction from context
type pgxTxKey struct{}

// Manager implements the transaction.Manager interface using pgx
type Manager struct {
	pool *pgxpool.Pool
}

// New creates a new Manager with the provided connection pool
func New(pool *pgxpool.Pool) *Manager {
	return &Manager{
		pool: pool,
	}
}

// Begin starts a new transaction
func (m *Manager) Begin(ctx context.Context) (context.Context, error) {
	tx, err := m.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin pgx transaction: %w", err)
	}

	txCtx := context.WithValue(ctx, pgxTxKey{}, tx)
	return txCtx, nil
}

// Commit commits the transaction
func (m *Manager) Commit(ctx context.Context) error {
	tx, err := getPgxTx(ctx)
	if err != nil {
		return fmt.Errorf("get transaction: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit pgx transaction: %w", err)
	}

	return nil
}

// Rollback aborts the transaction
func (m *Manager) Rollback(ctx context.Context) error {
	tx, err := getPgxTx(ctx)
	if err != nil {
		return fmt.Errorf("get transaction: %w", err)
	}

	if err := tx.Rollback(ctx); err != nil {
		return fmt.Errorf("rollback pgx transaction: %w", err)
	}

	return nil
}

// ExecTx executes a function within a transaction
func (m *Manager) ExecTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := m.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin pgx transaction: %w", err)
	}

	txCtx := context.WithValue(ctx, pgxTxKey{}, tx)

	if err := fn(txCtx); err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return fmt.Errorf("transaction failed: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit pgx transaction: %w", err)
	}

	return nil
}

// getPgxTx extracts the pgx.Tx from context
func getPgxTx(ctx context.Context) (pgx.Tx, error) {
	tx, ok := ctx.Value(pgxTxKey{}).(pgx.Tx)
	if !ok {
		return nil, errors.New("pgx transaction not found in context")
	}
	return tx, nil
}
