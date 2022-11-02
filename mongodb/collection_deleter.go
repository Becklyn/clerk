package mongodb

import (
	"context"
	"strings"

	"github.com/Becklyn/clerk/v2"

	"go.mongodb.org/mongo-driver/mongo"
)

type collectionDeleter struct {
	client   *mongo.Client
	database *clerk.Database
}

func newCollectionDeleter(connection *Connection, database *clerk.Database) *collectionDeleter {
	return &collectionDeleter{
		client:   connection.client,
		database: database,
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

	for i, name := range names {
		err := d.client.
			Database(d.database.Name).
			Collection(name).
			Drop(ctx)
		if err != nil {
			return i, err
		}
	}
	return len(names), nil
}
