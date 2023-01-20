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

	db, release, err := c.tryUseDB(ctx)
	defer release()
	if err != nil {
		return err
	}

	for _, data := range create.Data {
		stat := fmt.Sprintf("CREATE TABLE %s (data JSONB DEFAULT '{}' NOT NULL); CREATE UNIQUE INDEX %s_index on %s ((data->>'_id'));", data.Name, data.Name, data.Name)

		_, err := db.Client().Exec(createCtx, stat)
		if err != nil {
			return err
		}
	}
	return nil
}
