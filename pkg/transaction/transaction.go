package transaction

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// Transaction defines a common transaction interface
type Transaction interface {
	ExecTx(ctx context.Context, fn func(ctx context.Context) error) error
}

// TxManager is an implementation of transaction management
type TxManager struct {
	db *sql.DB
}

// NewTxManager creates a new TxManager instance
func NewTxManager(db *sql.DB) *TxManager {
	return &TxManager{
		db: db,
	}
}

// ExecTx executes a function within a transaction scope
func (tm *TxManager) ExecTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := tm.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	txCtx := context.WithValue(ctx, txKey{}, tx)

	if err := fn(txCtx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

// Key for retrieving transaction from context
type txKey struct{}

// GetTx retrieves transaction from context
func GetTx(ctx context.Context) (*sql.Tx, error) {
	tx, ok := ctx.Value(txKey{}).(*sql.Tx)
	if !ok {
		return nil, errors.New("transaction not found in context")
	}
	return tx, nil
}
