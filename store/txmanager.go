package store

import (
	"context"
	"fmt"

	"github.com/TakumaKurosawa/sqlc-common-transaction/pkg/db"
	"github.com/TakumaKurosawa/sqlc-common-transaction/pkg/transaction"
	"github.com/TakumaKurosawa/sqlc-common-transaction/store/poststore"
	"github.com/TakumaKurosawa/sqlc-common-transaction/store/poststore/postpgstore"
	"github.com/TakumaKurosawa/sqlc-common-transaction/store/userstore"
	"github.com/TakumaKurosawa/sqlc-common-transaction/store/userstore/userpgstore"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TxManager coordinates database transactions across multiple stores
type TxManager struct {
	pgxManager *transaction.PgxManager
	pool       *pgxpool.Pool
}

// NewTxManager creates a new TxManager
func NewTxManager(pool *pgxpool.Pool) *TxManager {
	return &TxManager{
		pgxManager: transaction.NewPgxManager(pool),
		pool:       pool,
	}
}

// ExecTx executes a function within a transaction context, providing access to user and post stores
func (m *TxManager) ExecTx(ctx context.Context, fn func(userStore userstore.Store, postStore poststore.Store) error) error {
	return m.pgxManager.ExecTx(ctx, func(ctx context.Context) error {
		tx, err := transaction.GetPgxTx(ctx)
		if err != nil {
			return fmt.Errorf("get transaction: %w", err)
		}

		q := db.New(tx)
		userStore := m.NewUserStore(q)
		postStore := m.NewPostStore(q)

		if err := fn(userStore, postStore); err != nil {
			return fmt.Errorf("transaction failed: %w", err)
		}

		return nil
	})
}

// NewUserStore creates a new user store with the given queries
func (m *TxManager) NewUserStore(q *db.Queries) userstore.Store {
	return userpgstore.New(q)
}

// NewPostStore creates a new post store with the given queries
func (m *TxManager) NewPostStore(q *db.Queries) poststore.Store {
	return postpgstore.New(q)
}
