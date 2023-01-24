package postgres

import (
	"context"
	"errors"
	"strings"
	"sync"

	"github.com/Becklyn/clerk/v4"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/lo"
)

var (
	ErrRelationDoesNotExist = errors.New("42P01")
)

type txCtxData struct{}

var txDataKey txCtxData

type transactor struct {
	conn *Connection
}

func newTransactor(conn *Connection) *transactor {
	return &transactor{
		conn: conn,
	}
}

func (t *transactor) ExecuteTransaction(ctx context.Context, fn clerk.TransactionFn) error {
	txCtx := ctx
	tx, isNested := newTransactionCtxFromCtx(ctx)
	if !isNested {
		txCtx, tx = t.bindNewTransactionCtx(ctx)
	}

	if err := fn(txCtx); err != nil {
		if isNested {
			return err
		}

		if err := tx.Rollback(ctx); err != nil {
			return err
		}

		if tableName, isNotExistsErr := t.isTableNotExistsError(err); isNotExistsErr {
			return t.createTableAndReRun(ctx, txCtx, tableName, fn)
		}
	}

	if !isNested {
		return tx.Commit(ctx)
	}
	return nil
}

func (t *transactor) bindNewTransactionCtx(ctx context.Context) (context.Context, *transactionCtx) {
	tx := newTransactionCtx()
	return context.WithValue(ctx, txDataKey, tx), tx
}

func (t *transactor) isTableNotExistsError(err error) (string, bool) {
	pgErr, ok := err.(*pgconn.PgError)
	if !ok || pgErr.SQLState() != ErrRelationDoesNotExist.Error() {
		return "", false
	}

	messageParts := strings.Split(pgErr.Message, "\"")
	if len(messageParts) != 3 {
		return "", false
	}
	return messageParts[1], true
}

func (t *transactor) createTableAndReRun(ctx context.Context, txCtx context.Context, tableName string, fn clerk.TransactionFn) error {
	tx, _ := newTransactionCtxFromCtx(txCtx)

	return t.ExecuteTransaction(ctx, func(nestedCtx context.Context) error {
		for _, dbName := range lo.Keys(tx.txs) {
			database := clerk.NewDatabase(dbName)

			if err := newCollectionCreator(t.conn, database).ExecuteCreate(ctx, &clerk.Create[*clerk.Collection]{
				Data: []*clerk.Collection{
					clerk.NewCollection(database, tableName),
				},
			}); err != nil {
				return err
			}
		}

		return fn(nestedCtx)
	})
}

type transactionCtx struct {
	sync.RWMutex
	txs map[string]pgx.Tx
}

func newTransactionCtx() *transactionCtx {
	return &transactionCtx{
		txs: map[string]pgx.Tx{},
	}
}

func newTransactionCtxFromCtx(ctx context.Context) (*transactionCtx, bool) {
	tx, ok := ctx.Value(txDataKey).(*transactionCtx)
	return tx, ok
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

func (t *transactionCtx) createOrUse(ctx context.Context, dbName string, pool *pgxpool.Pool) (*pgx.Conn, error) {
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
