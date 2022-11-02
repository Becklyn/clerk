package mongodb

import (
	"context"

	"github.com/Becklyn/clerk/v2"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type databaseQuerier struct {
	client *mongo.Client
}

func newDatabaseQuerier(connection *Connection) *databaseQuerier {
	return &databaseQuerier{
		client: connection.client,
	}
}

func (q *databaseQuerier) ExecuteQuery(
	ctx context.Context,
	query *clerk.Query[*clerk.Database],
) (<-chan *clerk.Database, error) {
	opts := options.ListDatabases()

	filters, err := resolveFilters(query.Filters)
	if err != nil {
		return nil, err
	}

	names, err := q.client.
		ListDatabaseNames(ctx, filters, opts)
	if err != nil {
		return nil, err
	}

	channel := make(chan *clerk.Database)

	go func() {
		defer close(channel)

		for _, name := range names {
			channel <- clerk.NewDatabase(name)
		}
	}()

	return channel, nil
}
