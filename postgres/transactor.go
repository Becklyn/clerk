package postgres

import (
	"context"
	"sync"

	"github.com/Becklyn/clerk/v3"
	"github.com/jackc/pgx/v5"
)

type transactor struct {
	dbCreator *databaseCreator
}

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
	conn *Connection
	txs  map[string]pgx.Tx
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

func (t *transactionCtx) useDb(ctx context.Context, db *DatabaseConnection) (*pgx.Conn, error) {
	t.Lock()
	defer t.Unlock()

	if _, ok := t.txs[db.name]; !ok {
		tx, err := db.client.Begin(ctx)
		if err != nil {
			return nil, err
		}

		t.txs[db.name] = tx
	}

	return t.txs[db.name].Conn(), nil
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

func getConn(ctx context.Context, conn *Connection, database *clerk.Database) (*pgx.Conn, func(), error) {
	db, release, err := tryUseDb(ctx, conn, database)
	if err != nil {
		return nil, release, err
	}

	if tx, ok := ctx.Value(txCtxData{}).(*transactionCtx); ok {
		pgConn, err := tx.useDb(ctx, db)
		return pgConn, release, err
	}

	return db.client, release, nil
}

func tryUseDb(ctx context.Context, conn *Connection, database *clerk.Database) (*DatabaseConnection, func(), error) {
	db, release, err := conn.UseDatabase(database.Name)
	if err != nil {
		if errCreate := newDatabaseCreator(conn).ExecuteCreate(ctx, &clerk.Create[*clerk.Database]{
			Data: []*clerk.Database{database},
		}); errCreate != nil {
			return nil, nil, err
		}

		return conn.UseDatabase(database.Name)
	}

	return db, release, err
}
