package mongodb

import (
	"context"
	"strings"

	"github.com/Becklyn/clerk/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type databaseDeleter struct {
	client *mongo.Client
}

func newDatabaseDeleter(connection *Connection) *databaseDeleter {
	return &databaseDeleter{
		client: connection.client,
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

	for i, name := range names {
		err := d.client.
			Database(name).
			Drop(ctx)
		if err != nil {
			return i, err
		}
	}
	return len(names), nil
}
