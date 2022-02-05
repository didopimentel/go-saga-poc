package persistence

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

// NewTxManager creates a new transaction manager based on connection pool.
// minConn is the number of connections alive at least.
// maxConn is the number of connections alive at most
func NewTxManager(ctx context.Context, addr string, minConn, maxConn int32) (*TxManager, error) {
	pgxConn, err := NewPool(ctx, addr, minConn, maxConn)
	if err != nil {
		return nil, err
	}

	return &TxManager{
		ConnPool: pgxConn,
	}, nil
}

func NewPool(ctx context.Context, addr string, minConn, maxConn int32) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(addr)
	if err != nil {
		return nil, err
	}
	// The defaults are located on top of pgxpool.pool.go
	config.MaxConns = maxConn
	config.MinConns = minConn

	pgxConn, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	return pgxConn, nil
}
