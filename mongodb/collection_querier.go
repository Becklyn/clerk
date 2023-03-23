package mongodb

import (
	"context"

	"github.com/Becklyn/clerk/v4"

	"go.mongodb.org/mongo-driver/mongo/options"
)

type collectionQuerier struct {
	connection *Connection
	database   *clerk.Database
}

func newCollectionQuerier(connection *Connection, database *clerk.Database) *collectionQuerier {
	return &collectionQuerier{
		connection: connection,
		database:   database,
	}
}

func (q *collectionQuerier) Count(
	ctx context.Context,
	query *clerk.Query[*clerk.Collection],
) (int64, error) {
	opts := options.ListCollections()

	filters, err := resolveFilters(query.Filters...)
	if err != nil {
		return 0, err
	}

	queryCtx, cancel := q.connection.config.GetContext(ctx)
	defer cancel()

	names, err := q.connection.client.
		Database(q.database.Name).
		ListCollectionNames(queryCtx, filters, opts)

	return int64(len(names)), err
}

func (q *collectionQuerier) ExecuteQuery(
	ctx context.Context,
	query *clerk.Query[*clerk.Collection],
) (<-chan *clerk.Collection, error) {
	opts := options.ListCollections()

	filters, err := resolveFilters(query.Filters...)
	if err != nil {
		return nil, err
	}

	queryCtx, cancel := q.connection.config.GetContext(ctx)
	defer cancel()

	names, err := q.connection.client.
		Database(q.database.Name).
		ListCollectionNames(queryCtx, filters, opts)
	if err != nil {
		return nil, err
	}

	channel := make(chan *clerk.Collection)

	go func() {
		defer close(channel)

		for _, name := range names {
			select {
			case <-queryCtx.Done():
				return
			default:
				channel <- clerk.NewCollection(q.database, name)
			}
		}
	}()

	return channel, nil
}
