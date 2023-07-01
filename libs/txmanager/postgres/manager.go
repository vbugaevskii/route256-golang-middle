package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgtype/pgxtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

var txKey = struct{}{}

type Manager struct {
	Pool *pgxpool.Pool
}

func NewTxManager(pool *pgxpool.Pool) *Manager {
	return &Manager{Pool: pool}
}

func (m *Manager) RunRepeatableRead(ctx context.Context, fn func(ctxTx context.Context) error) error {
	tx, err := m.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.RepeatableRead,
	})
	if err != nil {
		return fmt.Errorf("failed tx begin: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	ctxTx := context.WithValue(ctx, txKey, tx)
	if err = fn(ctxTx); err != nil {
		return fmt.Errorf("failed tx set value: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed tx commit: %w", err)
	}

	return nil
}

func (m Manager) GetQuerier(ctx context.Context) pgxtype.Querier {
	tx, ok := ctx.Value(txKey).(pgxtype.Querier)
	if ok {
		return tx
	}
	return m.Pool
}
