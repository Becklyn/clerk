package mongodb

import (
	"context"

	"github.com/Becklyn/clerk/v3"

	"go.mongodb.org/mongo-driver/mongo/options"
)

type databaseQuerier struct {
	connection *Connection
}

func newDatabaseQuerier(connection *Connection) *databaseQuerier {
	return &databaseQuerier{
		connection: connection,
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

	queryCtx, cancel := q.connection.config.GetContext(ctx)
	defer cancel()

	names, err := q.connection.client.
		ListDatabaseNames(queryCtx, filters, opts)
	if err != nil {
		return nil, err
	}

	channel := make(chan *clerk.Database)

	go func() {
		defer close(channel)

		for _, name := range names {
			select {
			case <-queryCtx.Done():
				return
			default:
				channel <- clerk.NewDatabase(name)
			}
		}
	}()

	return channel, nil
}
