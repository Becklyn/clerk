package postgres

import (
	"context"
	"sync"

	"github.com/Becklyn/clerk/v4"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type txCtxData struct{}

var txDataKey txCtxData

type transactor struct{}

func newTransactor() *transactor {
	return &transactor{}
}

func (t *transactor) ExecuteTransaction(ctx context.Context, fn clerk.TransactionFn) error {
	txCtx, isNested, tx := useTx(ctx)

	if err := fn(txCtx); err != nil {
		if !isNested {
			if err := tx.Rollback(ctx); err != nil {
				return err
			}
		}
		return err
	}

	if !isNested {
		return tx.Commit(ctx)
	}
	return nil
}

type transactionCtx struct {
	sync.RWMutex
	txs map[string]pgx.Tx
}

func (t *transactionCtx) Rollback(ctx context.Context) error {
	t.RLock()
	defer t.RUnlock()

	for _, tx := range t.txs {
		if err := tx.Rollback(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (t *transactionCtx) Commit(ctx context.Context) error {
	t.RLock()
	defer t.RUnlock()

	for _, tx := range t.txs {
		if err := tx.Commit(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (t *transactionCtx) useDb(ctx context.Context, dbName string, pool *pgxpool.Pool) (*pgx.Conn, error) {
	t.Lock()
	defer t.Unlock()

	if _, ok := t.txs[dbName]; !ok {
		tx, err := pool.BeginTx(ctx, pgx.TxOptions{
			AccessMode: pgx.ReadWrite,
			IsoLevel:   pgx.ReadUncommitted,
		})
		if err != nil {
			return nil, err
		}

		t.txs[dbName] = tx
	}

	return t.txs[dbName].Conn(), nil
}

func useTx(ctx context.Context) (context.Context, bool, *transactionCtx) {
	if tx, ok := ctx.Value(txDataKey).(*transactionCtx); ok {
		return ctx, true, tx
	}

	tx := &transactionCtx{
		txs: map[string]pgx.Tx{},
	}
	return context.WithValue(ctx, txDataKey, tx), false, tx
}
