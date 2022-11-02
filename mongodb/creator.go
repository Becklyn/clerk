package mongodb

import (
	"context"

	"github.com/Becklyn/clerk/v2"

	"go.mongodb.org/mongo-driver/mongo"
)

type creator[T any] struct {
	client     *mongo.Client
	collection *clerk.Collection
}

func newCreator[T any](connection *Connection, collection *clerk.Collection) *creator[T] {
	return &creator[T]{
		client:     connection.client,
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

	_, err := c.client.
		Database(c.collection.Database.Name).
		Collection(c.collection.Name).
		InsertMany(ctx, data)

	return err
}
