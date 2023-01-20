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
	dbCons map[string]*clerk.UsagePool[*DatabaseConnection]
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
		dbCons: map[string]*clerk.UsagePool[*DatabaseConnection]{},
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
					_ = usagePool.Get().Close()
					delete(c.dbCons, database)
				}
			}
		}
	}
}

// TODO: DONT EXPORT THIS
func (c *Connection) Client() *pgx.Conn {
	return c.client
}

func (c *Connection) UseDatabase(database string) (*DatabaseConnection, func(), error) {
	usagePool, ok := c.dbCons[database]
	if !ok {
		dbCon, err := NewDatabaseConnection(c, database)
		if err != nil {
			return nil, nil, err
		}

		usagePool = clerk.NewUsagePool(dbCon, c.config.Timeout)
		c.dbCons[database] = usagePool
	}
	return usagePool.Get(), usagePool.Release, nil
}

func (c *Connection) Close(handler func(err error)) {
	for _, usagePool := range c.dbCons {
		err := usagePool.Get().Close()
		if err != nil && handler != nil {
			handler(err)
		}
	}
	if err := c.client.Close(c.ctx); err != nil && handler != nil {
		handler(err)
	}
}
