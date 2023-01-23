package postgres

import (
	"context"
	"fmt"
	"sync"

	"github.com/Becklyn/clerk/v4"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Connection struct {
	ctx    context.Context
	config Config
	pool   *pgxpool.Pool

	sync.Mutex
	dbPools map[string]*pgxpool.Pool
}

func NewConnection(
	ctx context.Context,
	config Config,
) (*Connection, error) {
	pool, err := pgxpool.New(ctx, string(config.Host))
	if err != nil {
		return nil, err
	}

	pingCtx, pingCancel := config.GetContext(ctx)
	defer pingCancel()
	if err = pool.Ping(pingCtx); err != nil {
		return nil, err
	}

	return &Connection{
		ctx:     ctx,
		config:  config,
		pool:    pool,
		dbPools: map[string]*pgxpool.Pool{},
	}, nil
}

func (c *Connection) useDatabase(ctx context.Context, database string) (*pgx.Conn, func(), error) {
	pool, err := c.tryUseDb(ctx, database)
	if err != nil {
		return nil, func() {}, err
	}

	if tx, ok := ctx.Value(txCtxData{}).(*transactionCtx); ok {
		pgConn, err := tx.useDb(ctx, database, pool)
		return pgConn, func() {}, err
	}

	conn, err := pool.Acquire(ctx)
	if err != nil {
		return nil, func() {}, err
	}

	return conn.Conn(), conn.Release, nil
}

func (c *Connection) Close() {
	c.Lock()
	defer c.Unlock()

	for _, pool := range c.dbPools {
		pool.Close()
	}
	c.pool.Close()
}

func (c *Connection) getDbPool(database string) (*pgxpool.Pool, error) {
	c.Lock()
	defer c.Unlock()

	dbPool, ok := c.dbPools[database]
	if ok {
		return dbPool, nil
	}

	dbPool, err := pgxpool.New(c.ctx, fmt.Sprintf("%s/%s", c.config.Host, database))
	if err != nil {
		return nil, err
	}

	c.dbPools[database] = dbPool

	return dbPool, nil
}

func (c *Connection) tryUseDb(ctx context.Context, database string) (*pgxpool.Pool, error) {
	pool, err := c.getDbPool(database)
	if err != nil {
		if errCreate := newDatabaseCreator(c).ExecuteCreate(ctx, &clerk.Create[*clerk.Database]{
			Data: []*clerk.Database{clerk.NewDatabase(database)},
		}); errCreate != nil {
			return nil, err
		}

		return c.getDbPool(database)
	}

	return pool, err
}
