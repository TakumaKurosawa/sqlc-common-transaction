package pgxtransaction

import (
	"context"
	"errors"
	"fmt"

	"github.com/TakumaKurosawa/sqlc-common-transaction/pkg/db"
	"github.com/TakumaKurosawa/sqlc-common-transaction/store/poststore"
	"github.com/TakumaKurosawa/sqlc-common-transaction/store/poststore/postpgstore"
	"github.com/TakumaKurosawa/sqlc-common-transaction/store/userstore"
	"github.com/TakumaKurosawa/sqlc-common-transaction/store/userstore/userpgstore"
	"github.com/TakumaKurosawa/sqlc-common-transaction/transaction"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// pgxTxKey is a key for retrieving transaction from context
type pgxTxKey struct{}

// Manager implements transaction management using pgx
type Manager struct {
	pool *pgxpool.Pool
}

// New creates a new Manager
func New(pool *pgxpool.Pool) transaction.Manager {
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
func (m *Manager) ExecTx(ctx context.Context, fn func(userStore userstore.Store, postStore poststore.Store) error) error {
	tx, err := m.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin pgx transaction: %w", err)
	}

	// The transaction context is created but unused to match interface expectations
	_ = context.WithValue(ctx, pgxTxKey{}, tx)

	q := db.New(tx)
	userStore := m.newUserStore(q)
	postStore := m.newPostStore(q)

	if err := fn(userStore, postStore); err != nil {
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

// WithTx executes a function using an existing transaction context
func (m *Manager) WithTx(ctx context.Context, fn func(userStore userstore.Store, postStore poststore.Store) error) error {
	tx, err := getPgxTx(ctx)
	if err != nil {
		return fmt.Errorf("get transaction: %w", err)
	}

	q := db.New(tx)
	userStore := m.newUserStore(q)
	postStore := m.newPostStore(q)

	if err := fn(userStore, postStore); err != nil {
		return fmt.Errorf("transaction operation failed: %w", err)
	}

	return nil
}

// getPgxTx retrieves transaction from context
func getPgxTx(ctx context.Context) (pgx.Tx, error) {
	tx, ok := ctx.Value(pgxTxKey{}).(pgx.Tx)
	if !ok {
		return nil, errors.New("pgx transaction not found in context")
	}
	return tx, nil
}

// newUserStore creates a new user store
func (m *Manager) newUserStore(q *db.Queries) userstore.Store {
	return userpgstore.New(q)
}

// newPostStore creates a new post store
func (m *Manager) newPostStore(q *db.Queries) poststore.Store {
	return postpgstore.New(q)
}
