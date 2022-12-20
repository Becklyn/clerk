package mongodb

import (
	"context"
	"strings"

	"github.com/Becklyn/clerk/v3"
)

type databaseDeleter struct {
	connection *Connection
}

func newDatabaseDeleter(connection *Connection) *databaseDeleter {
	return &databaseDeleter{
		connection: connection,
	}
}

func (d *databaseDeleter) ExecuteDelete(
	ctx context.Context,
	delete *clerk.Delete[*clerk.Database],
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
			Database(name).
			Drop(deleteCtx)
		if err != nil {
			return i, err
		}
	}
	return len(names), nil
}
