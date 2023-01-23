package postgres

import (
	"context"
	"fmt"

	"github.com/Becklyn/clerk/v3"
	"github.com/jackc/pgx/v5"
)

type Connection struct {
	ctx    context.Context
	config Config
	client *pgx.Conn
	dbCons map[string]*clerk.UsagePool[*pgx.Conn]
}

func NewConnection(
	ctx context.Context,
	config Config,
) (*Connection, error) {
	client, err := pgx.Connect(ctx, string(config.Host))
	if err != nil {
		return nil, err
	}

	pingCtx, pingCancel := config.GetContext(ctx)
	defer pingCancel()
	if err = client.Ping(pingCtx); err != nil {
		return nil, err
	}

	c := &Connection{
		ctx:    ctx,
		config: config,
		client: client,
		dbCons: map[string]*clerk.UsagePool[*pgx.Conn]{},
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
			for database, usagePool := range c.dbCons {
				if usagePool.IsUnused() {
					fmt.Printf("Closing unused database conn to %s\n", database)
					_ = usagePool.Get().Close(c.ctx)
					delete(c.dbCons, database)
				}
			}
		}
	}
}

func (c *Connection) useDatabase(ctx context.Context, database string) (*pgx.Conn, func(), error) {
	db, release, err := c.tryUseDb(ctx, database)
	if err != nil {
		return nil, release, err
	}

	if tx, ok := ctx.Value(txCtxData{}).(*transactionCtx); ok {
		pgConn, err := tx.useDb(ctx, database, db)
		return pgConn, release, err
	}

	return db, release, nil
}

func (c *Connection) Close(handler func(err error)) {
	for _, usagePool := range c.dbCons {
		err := usagePool.Get().Close(c.ctx)
		if err != nil && handler != nil {
			handler(err)
		}
	}
	if err := c.client.Close(c.ctx); err != nil && handler != nil {
		handler(err)
	}
}

func (c *Connection) getDbClient(database string) (*pgx.Conn, func(), error) {
	usagePool, ok := c.dbCons[database]
	if !ok {
		client, err := pgx.Connect(c.ctx, fmt.Sprintf("%s/%s", c.config.Host, database))
		if err != nil {
			return nil, nil, err
		}

		usagePool = clerk.NewUsagePool(client, c.config.Timeout)
		c.dbCons[database] = usagePool
	}
	return usagePool.Get(), usagePool.Release, nil
}

func (c *Connection) tryUseDb(ctx context.Context, database string) (*pgx.Conn, func(), error) {
	db, release, err := c.getDbClient(database)
	if err != nil {
		if errCreate := newDatabaseCreator(c).ExecuteCreate(ctx, &clerk.Create[*clerk.Database]{
			Data: []*clerk.Database{clerk.NewDatabase(database)},
		}); errCreate != nil {
			return nil, nil, err
		}

		return c.getDbClient(database)
	}

	return db, release, err
}
