package mongodb

import (
	"context"

	"github.com/Becklyn/clerk/v2"

	"go.mongodb.org/mongo-driver/mongo"
)

type deleter[T any] struct {
	client     *mongo.Client
	collection *clerk.Collection
}

func newDeleter[T any](connection *Connection, collection *clerk.Collection) *deleter[T] {
	return &deleter[T]{
		client:     connection.client,
		collection: collection,
	}
}

func (d *deleter[T]) ExecuteDelete(
	ctx context.Context,
	delete *clerk.Delete[T],
) (int, error) {
	filters, err := resolveFilters(delete.Filters)
	if err != nil {
		return 0, err
	}

	result, err := d.client.
		Database(d.collection.Database.Name).
		Collection(d.collection.Name).
		DeleteMany(ctx, filters)

	return int(result.DeletedCount), err
}
