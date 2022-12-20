package mongodb

import (
	"context"

	"github.com/Becklyn/clerk/v3"

	"go.mongodb.org/mongo-driver/mongo/options"
)

type updater[T any] struct {
	connection *Connection
	collection *clerk.Collection
}

func newUpdater[T any](connection *Connection, collection *clerk.Collection) *updater[T] {
	return &updater[T]{
		connection: connection,
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

	updaterCtx, cancel := u.connection.config.GetContext(ctx)
	defer cancel()

	_, err = u.connection.client.
		Database(u.collection.Database.Name).
		Collection(u.collection.Name).
		ReplaceOne(updaterCtx, filters, update.Data, opts)

	return err
}
