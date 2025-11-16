package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type txKey struct{}

type PG struct {
	Pool *pgxpool.Pool
}

func New(conn string) (*PG, error) {
	pool, err := pgxpool.New(context.Background(), conn)
	if err != nil {
		return nil, err
	}
	return &PG{Pool: pool}, nil
}
func (pg *PG) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := pg.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	ctxWithTx := context.WithValue(ctx, txKey{}, tx)

	if err := fn(ctxWithTx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func TxFromContext(ctx context.Context) (pgx.Tx, bool) {
	v := ctx.Value(txKey{})
	if v == nil {
		return nil, false
	}
	tx, ok := v.(pgx.Tx)
	return tx, ok
}
