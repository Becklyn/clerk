package mongodb

import (
	"context"

	"github.com/Becklyn/clerk/v3"
)

type deleter[T any] struct {
	connection *Connection
	collection *clerk.Collection
}

func newDeleter[T any](connection *Connection, collection *clerk.Collection) *deleter[T] {
	return &deleter[T]{
		connection: connection,
		collection: collection,
	}
}

func (d *deleter[T]) ExecuteDelete(
	ctx context.Context,
	delete *clerk.Delete[T],
) (int, error) {
	filters, err := resolveFilters(delete.Filters...)
	if err != nil {
		return 0, err
	}

	deleteCtx, cancel := d.connection.config.GetContext(ctx)
	defer cancel()

	result, err := d.connection.client.
		Database(d.collection.Database.Name).
		Collection(d.collection.Name).
		DeleteMany(deleteCtx, filters)

	return int(result.DeletedCount), err
}
