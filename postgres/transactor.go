package postgres

import (
	"context"
	"sync"

	"github.com/Becklyn/clerk/v3"
	"github.com/jackc/pgx/v5"
)

type transactor struct{}

func newTransactor() *transactor {
	return &transactor{}
}

func (t *transactor) ExecuteTransaction(ctx context.Context, fn clerk.TransactionFn) error {
	txCtx, tx := useTx(ctx)
	defer tx.Rollback(ctx)

	if err := fn(txCtx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

type transactionCtx struct {
	sync.RWMutex
	txs map[string]pgx.Tx
}

func (t *transactionCtx) Rollback(ctx context.Context) {
	t.RLock()
	defer t.RUnlock()

	for _, tx := range t.txs {
		_ = tx.Rollback(ctx)
	}
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

func (t *transactionCtx) useDb(ctx context.Context, dbName string, db *pgx.Conn) (*pgx.Conn, error) {
	t.Lock()
	defer t.Unlock()

	if _, ok := t.txs[dbName]; !ok {
		tx, err := db.Begin(ctx)
		if err != nil {
			return nil, err
		}

		t.txs[dbName] = tx
	}

	return t.txs[dbName].Conn(), nil
}

type txCtxData struct{}

func useTx(ctx context.Context) (context.Context, *transactionCtx) {
	if tx, ok := ctx.Value(txCtxData{}).(*transactionCtx); ok {
		return ctx, tx
	}

	tx := &transactionCtx{
		txs: map[string]pgx.Tx{},
	}
	return context.WithValue(ctx, txCtxData{}, tx), tx
}
