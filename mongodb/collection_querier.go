package mongodb

import (
	"context"

	"github.com/Becklyn/clerk/v2"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type collectionQuerier struct {
	client   *mongo.Client
	database *clerk.Database
}

func newCollectionQuerier(connection *Connection, database *clerk.Database) *collectionQuerier {
	return &collectionQuerier{
		client:   connection.client,
		database: database,
	}
}

func (q *collectionQuerier) ExecuteQuery(
	ctx context.Context,
	query *clerk.Query[*clerk.Collection],
) (<-chan *clerk.Collection, error) {
	opts := options.ListCollections()

	filters, err := resolveFilters(query.Filters)
	if err != nil {
		return nil, err
	}

	names, err := q.client.
		Database(q.database.Name).
		ListCollectionNames(ctx, filters, opts)
	if err != nil {
		return nil, err
	}

	channel := make(chan *clerk.Collection)

	go func() {
		defer close(channel)

		for _, name := range names {
			channel <- clerk.NewCollection(q.database, name)
		}
	}()

	return channel, nil
}
