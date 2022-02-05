package persistence

import (
	"context"
	"errors"
	"fmt"
	"github.com/didopimentel/go-saga-poc/domain"
	"github.com/jackc/pgx/v4"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/jackc/pgconn"
)

type txCtxKey string

const txKey = txCtxKey("transaction")

type TxManager struct {
	ConnPool *pgxpool.Pool
}

func (b TxManager) WithTx(ctx context.Context, f domain.TransactionFunc) error {
	if domain.InTX(ctx) {
		return &domain.TransactionError{Cause: errors.New("already in transaction")}
	}
	var err error
	var tx pgx.Tx
	tx, err = b.ConnPool.Begin(ctx)
	if err != nil {
		return &domain.TransactionError{Cause: fmt.Errorf("cannot begin a transaction: %w", err)}
	}
	ctxWithTx := context.WithValue(domain.ContextWithTx(ctx), txKey, tx)

	defer func() {
		if p := recover(); p != nil {
			// ensure a rollback attempt and panic again
			_ = tx.Rollback(ctx) //nolint:errcheck
			panic(p)
		}
	}()

	if err = f(ctxWithTx); err != nil {
		if rollBackErr := tx.Rollback(ctx); rollBackErr != nil {
			return rollBackErr
		}

		return err
	}

	return tx.Commit(ctx)
}

// querier should be used when either a transaction or a common connection could be used
type querier interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

func (b TxManager) querier(ctx context.Context) querier {
	if btx, ok := ctx.Value(txKey).(*pgxpool.Tx); ok {
		return btx
	}

	return b.ConnPool
}

func (b TxManager) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	return b.querier(ctx).Exec(ctx, query, args...)
}

func (b TxManager) Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	return b.querier(ctx).Query(ctx, query, args...)
}

func (b TxManager) QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	return b.querier(ctx).QueryRow(ctx, query, args...)
}
