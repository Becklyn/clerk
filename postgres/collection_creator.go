package postgres

import (
	"context"
	"fmt"

	"github.com/Becklyn/clerk/v3"
)

type collectionCreator struct {
	collectionBase
}

func newCollectionCreator(conn *Connection, database *clerk.Database) *collectionCreator {
	return &collectionCreator{
		*newCollectionBase(conn, database),
	}
}

func (c *collectionCreator) ExecuteCreate(
	ctx context.Context,
	create *clerk.Create[*clerk.Collection],
) error {
	createCtx, cancel := c.conn.config.GetContext(ctx)
	defer cancel()

	dbConn, release, err := c.conn.useDatabase(createCtx, c.database.Name)
	defer release()
	if err != nil {
		return err
	}

	for _, data := range create.Data {
		stat := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (data JSONB DEFAULT '{}' NOT NULL); CREATE UNIQUE INDEX IF NOT EXISTS %s_pk on %s ((data->>'_id'));", data.Name, data.Name, data.Name)

		_, err := dbConn.Exec(createCtx, stat)
		if err != nil {
			return err
		}
	}
	return nil
}
