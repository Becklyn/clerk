package mongodb

import (
	"context"
	"strings"

	"github.com/Becklyn/clerk/v2"

	"go.mongodb.org/mongo-driver/mongo"
)

type indexDeleter struct {
	client     *mongo.Client
	collection *clerk.Collection
}

func newIndexDeleter(connection *Connection, collection *clerk.Collection) *indexDeleter {
	return &indexDeleter{
		client:     connection.client,
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

	if len(names) == 0 {
		_, err := d.client.
			Database(d.collection.Database.Name).
			Collection(d.collection.Name).
			Indexes().
			DropAll(ctx)

		return 0, err
	}

	for i, name := range names {
		_, err := d.client.
			Database(d.collection.Database.Name).
			Collection(d.collection.Name).
			Indexes().
			DropOne(ctx, name)
		if err != nil {
			return i, err
		}
	}
	return len(names), nil
}
