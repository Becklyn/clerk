package mongodb

import (
	"context"
	"strings"

	"github.com/Becklyn/clerk/v3"
)

type collectionDeleter struct {
	connection *Connection
	database   *clerk.Database
}

func newCollectionDeleter(connection *Connection, database *clerk.Database) *collectionDeleter {
	return &collectionDeleter{
		connection: connection,
		database:   database,
	}
}

func (d *collectionDeleter) ExecuteDelete(
	ctx context.Context,
	delete *clerk.Delete[*clerk.Collection],
) (int, error) {
	names := []string{}
	for _, filter := range delete.Filters {
		switch filter.(type) {
		case *clerk.Equals:
			if strings.ToLower(filter.Key()) == "name" {
				names = append(names, filter.Value().(string))
			}
		}
	}

	deleteCtx, cancel := d.connection.config.GetContext(ctx)
	defer cancel()

	for i, name := range names {
		err := d.connection.client.
			Database(d.database.Name).
			Collection(name).
			Drop(deleteCtx)
		if err != nil {
			return i, err
		}
	}
	return len(names), nil
}
