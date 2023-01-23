package mongodb

import (
	"context"

	"github.com/Becklyn/clerk/v4"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type querier[T any] struct {
	connection *Connection
	collection *clerk.Collection
}

func newQuerier[T any](connection *Connection, collection *clerk.Collection) *querier[T] {
	return &querier[T]{
		connection: connection,
		collection: collection,
	}
}

func (q *querier[T]) ExecuteQuery(
	ctx context.Context,
	query *clerk.Query[T],
) (<-chan T, error) {
	opts := options.Find()

	if query.Range != nil {
		opts.SetSkip(int64(query.Range.SkipValue))
		opts.SetLimit(int64(query.Range.TakeValue))
	}

	if len(query.Sorting) > 0 {
		sort := bson.D{}
		for key, order := range query.Sorting {
			if order.IsAscending {
				sort = append(sort, bson.E{
					Key:   key,
					Value: 1,
				})
			} else {
				sort = append(sort, bson.E{
					Key:   key,
					Value: -1,
				})
			}
		}
		opts.SetSort(sort)
	}

	filters, err := resolveFilters(query.Filters...)
	if err != nil {
		return nil, err
	}

	queryCtx, cancel := q.connection.config.GetContext(ctx)
	defer cancel()

	cursor, err := q.connection.client.
		Database(q.collection.Database.Name).
		Collection(q.collection.Name).
		Find(queryCtx, filters, opts)
	if err != nil {
		return nil, err
	}

	channel := make(chan T)

	go func() {
		defer cursor.Close(queryCtx)
		defer close(channel)

		for cursor.Next(queryCtx) {
			var result T
			if err := cursor.Decode(&result); err != nil {
				return
			}
			channel <- result
		}
	}()

	return channel, nil
}
