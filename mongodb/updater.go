package mongodb

import (
	"context"

	"github.com/Becklyn/clerk/v2"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type updater[T any] struct {
	client     *mongo.Client
	collection *clerk.Collection
}

func newUpdater[T any](connection *Connection, collection *clerk.Collection) *updater[T] {
	return &updater[T]{
		client:     connection.client,
		collection: collection,
	}
}

func (u *updater[T]) ExecuteUpdate(ctx context.Context, update *clerk.Update[T]) error {
	opts := options.Replace().
		SetUpsert(update.ShouldUpsert)

	filters, err := resolveFilters(update.Filters)
	if err != nil {
		return err
	}

	_, err = u.client.
		Database(u.collection.Database.Name).
		Collection(u.collection.Name).
		ReplaceOne(ctx, filters, update.Data, opts)

	return err
}
