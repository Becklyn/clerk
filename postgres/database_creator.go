package postgres

import (
	"context"
	"github.com/Becklyn/clerk/v3"
)

type databaseCreator struct {
	conn *Connection
}

func newDatabaseCreator(conn *Connection) *databaseCreator {
	return &databaseCreator{
		conn: conn,
	}
}

func (c *databaseCreator) ExecuteCreate(
	ctx context.Context,
	create *clerk.Create[*clerk.Database],
) error {
	createCtx, cancel := c.conn.config.GetContext(ctx)
	defer cancel()

	for _, data := range create.Data {
		stat := "CREATE DATABASE " + data.Name
		if _, err := c.conn.client.Exec(createCtx, stat); err != nil {
			return err
		}
	}
	return nil
}
