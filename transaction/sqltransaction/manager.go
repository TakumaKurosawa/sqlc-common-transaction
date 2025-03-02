package sqltransaction

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/TakumaKurosawa/sqlc-common-transaction/model"
	"github.com/TakumaKurosawa/sqlc-common-transaction/store/poststore"
	"github.com/TakumaKurosawa/sqlc-common-transaction/store/userstore"
	"github.com/TakumaKurosawa/sqlc-common-transaction/transaction"
	"github.com/google/uuid"
)

// txKey is a key for retrieving transaction from context
type txKey struct{}

// Manager implements transaction management using standard SQL
type Manager struct {
	db *sql.DB
}

// New creates a new Manager
func New(db *sql.DB) transaction.Manager {
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
func (m *Manager) ExecTx(ctx context.Context, fn func(userStore userstore.Store, postStore poststore.Store) error) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	// In a real application, we would need to create actual SQL store implementations here.
	// Since the sample code doesn't have standard SQL store implementations, we're using dummy implementations.
	userStore := &dummyUserStore{}
	postStore := &dummyPostStore{}

	if err := fn(userStore, postStore); err != nil {
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

// WithTx executes a function using an existing transaction context
func (m *Manager) WithTx(ctx context.Context, fn func(userStore userstore.Store, postStore poststore.Store) error) error {
	// Only retrieve the transaction from context for validation.
	// In this sample, we'll use dummy stores.
	if _, err := getTx(ctx); err != nil {
		return fmt.Errorf("get transaction: %w", err)
	}

	// In a real application, we would need to create actual SQL store implementations here.
	// Since the sample code doesn't have standard SQL store implementations, we're using dummy implementations.
	userStore := &dummyUserStore{}
	postStore := &dummyPostStore{}

	if err := fn(userStore, postStore); err != nil {
		return fmt.Errorf("transaction operation failed: %w", err)
	}

	return nil
}

// getTx retrieves transaction from context
func getTx(ctx context.Context) (*sql.Tx, error) {
	tx, ok := ctx.Value(txKey{}).(*sql.Tx)
	if !ok {
		return nil, fmt.Errorf("transaction not found in context")
	}
	return tx, nil
}

// dummyUserStore and dummyPostStore are sample implementations for the interface.
// In a real application, these would be replaced with actual SQL implementations.

type dummyUserStore struct{}

func (d *dummyUserStore) CreateUser(ctx context.Context, name, email string) (model.User, error) {
	return model.User{}, fmt.Errorf("not implemented")
}

func (d *dummyUserStore) GetUser(ctx context.Context, id uuid.UUID) (model.User, error) {
	return model.User{}, fmt.Errorf("not implemented")
}

func (d *dummyUserStore) ListUsers(ctx context.Context) ([]model.User, error) {
	return nil, fmt.Errorf("not implemented")
}

func (d *dummyUserStore) UpdateUser(ctx context.Context, id uuid.UUID, name, email string) (model.User, error) {
	return model.User{}, fmt.Errorf("not implemented")
}

func (d *dummyUserStore) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return fmt.Errorf("not implemented")
}

type dummyPostStore struct{}

func (d *dummyPostStore) CreatePost(ctx context.Context, userID uuid.UUID, title, content string) (model.Post, error) {
	return model.Post{}, fmt.Errorf("not implemented")
}

func (d *dummyPostStore) GetPost(ctx context.Context, id uuid.UUID) (model.Post, error) {
	return model.Post{}, fmt.Errorf("not implemented")
}

func (d *dummyPostStore) ListPostsByUser(ctx context.Context, userID uuid.UUID) ([]model.Post, error) {
	return nil, fmt.Errorf("not implemented")
}

func (d *dummyPostStore) UpdatePost(ctx context.Context, id uuid.UUID, title, content string) (model.Post, error) {
	return model.Post{}, fmt.Errorf("not implemented")
}

func (d *dummyPostStore) DeletePost(ctx context.Context, id uuid.UUID) error {
	return fmt.Errorf("not implemented")
}
