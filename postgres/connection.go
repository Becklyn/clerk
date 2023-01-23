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
	dbPools map[string]*clerk.UsagePool[*pgxpool.Pool]
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

	c := &Connection{
		ctx:     ctx,
		config:  config,
		pool:    pool,
		dbPools: map[string]*clerk.UsagePool[*pgxpool.Pool]{},
	}
	go c.dbConsCleanupTask()
	return c, nil
}

func (c *Connection) dbConsCleanupTask() {
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			c.Lock()
			for database, usagePool := range c.dbPools {
				if usagePool.IsUnused() {
					fmt.Printf("Closing unused database conn to %s\n", database)
					usagePool.Get().Close()
					delete(c.dbPools, database)
				}
			}
			c.Unlock()
		}
	}
}

func (c *Connection) useDatabase(ctx context.Context, database string) (*pgx.Conn, func(), error) {
	pool, release, err := c.tryUseDb(ctx, database)
	if err != nil {
		return nil, release, err
	}

	if tx, ok := ctx.Value(txCtxData{}).(*transactionCtx); ok {
		pgConn, err := tx.useDb(ctx, database, pool)
		return pgConn, release, err
	}

	conn, err := pool.Acquire(ctx)
	release = func() {
		conn.Release()
		release()
	}

	if err != nil {
		return nil, release, err
	}

	return conn.Conn(), release, nil
}

func (c *Connection) Close() {
	c.Lock()
	defer c.Unlock()

	for _, usagePool := range c.dbPools {
		usagePool.Get().Close()
	}
	c.pool.Close()
}

func (c *Connection) getDbPool(database string) (*pgxpool.Pool, func(), error) {
	c.Lock()
	defer c.Unlock()

	usagePool, ok := c.dbPools[database]
	if !ok {
		dbPool, err := pgxpool.New(c.ctx, fmt.Sprintf("%s/%s", c.config.Host, database))
		if err != nil {
			return nil, nil, err
		}

		usagePool = clerk.NewUsagePool(dbPool, c.config.Timeout)
		c.dbPools[database] = usagePool
	}

	return usagePool.Get(), usagePool.Release, nil
}

func (c *Connection) tryUseDb(ctx context.Context, database string) (*pgxpool.Pool, func(), error) {
	pool, release, err := c.getDbPool(database)
	if err != nil {
		if errCreate := newDatabaseCreator(c).ExecuteCreate(ctx, &clerk.Create[*clerk.Database]{
			Data: []*clerk.Database{clerk.NewDatabase(database)},
		}); errCreate != nil {
			return nil, nil, err
		}

		return c.getDbPool(database)
	}

	return pool, release, err
}
