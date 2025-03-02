package transaction

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PgxTxKey is a transaction key for pgx
type pgxTxKey struct{}

// PgxManager is an implementation of transaction management using pgx
type PgxManager struct {
	pool *pgxpool.Pool
}

// NewPgxManager creates a new PgxManager instance
func NewPgxManager(pool *pgxpool.Pool) *PgxManager {
	return &PgxManager{
		pool: pool,
	}
}

// ExecTx executes a function within a pgx transaction scope
func (p *PgxManager) ExecTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin pgx transaction: %w", err)
	}

	txCtx := context.WithValue(ctx, pgxTxKey{}, tx)

	err = fn(txCtx)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit pgx transaction: %w", err)
	}

	return nil
}

// GetPgxTx retrieves pgx transaction from context
func GetPgxTx(ctx context.Context) (pgx.Tx, error) {
	tx, ok := ctx.Value(pgxTxKey{}).(pgx.Tx)
	if !ok {
		return nil, errors.New("pgx transaction not found in context")
	}
	return tx, nil
}
