package mongodb

import (
	"context"

	"github.com/Becklyn/clerk/v4"
)

type creator[T any] struct {
	connection *Connection
	collection *clerk.Collection
}

func newCreator[T any](connection *Connection, collection *clerk.Collection) *creator[T] {
	return &creator[T]{
		connection: connection,
		collection: collection,
	}
}

func (c *creator[T]) ExecuteCreate(
	ctx context.Context,
	create *clerk.Create[T],
) error {
	data := make([]any, len(create.Data))
	for i, item := range create.Data {
		data[i] = item
	}

	createCtx, cancel := c.connection.config.GetContext(ctx)
	defer cancel()

	_, err := c.connection.client.
		Database(c.collection.Database.Name).
		Collection(c.collection.Name).
		InsertMany(createCtx, data)

	return err
}
