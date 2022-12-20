package mongodb

import (
	"context"
	"strings"

	"github.com/Becklyn/clerk/v3"
)

type indexDeleter struct {
	connection *Connection
	collection *clerk.Collection
}

func newIndexDeleter(connection *Connection, collection *clerk.Collection) *indexDeleter {
	return &indexDeleter{
		connection: connection,
		collection: collection,
	}
}

func (d *indexDeleter) ExecuteDelete(
	ctx context.Context,
	delete *clerk.Delete[*clerk.Index],
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

	if len(names) == 0 {
		_, err := d.connection.client.
			Database(d.collection.Database.Name).
			Collection(d.collection.Name).
			Indexes().
			DropAll(deleteCtx)

		return 0, err
	}

	for i, name := range names {
		_, err := d.connection.client.
			Database(d.collection.Database.Name).
			Collection(d.collection.Name).
			Indexes().
			DropOne(deleteCtx, name)
		if err != nil {
			return i, err
		}
	}
	return len(names), nil
}
