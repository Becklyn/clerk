package postgres

import (
	"context"

	"github.com/Becklyn/clerk/v3"
	"github.com/jackc/pgx/v5"
	"github.com/xdg-go/jibby"
	"go.mongodb.org/mongo-driver/bson"
)

type querier[T any] struct {
	conn       *Connection
	collection *clerk.Collection
}

func newQuerier[T any](conn *Connection, collection *clerk.Collection) *querier[T] {
	return &querier[T]{
		conn:       conn,
		collection: collection,
	}
}

func (q *querier[T]) ExecuteQuery(
	ctx context.Context,
	query *clerk.Query[T],
) (<-chan T, error) {
	statBuilder := statementBuilder().
		Select("(data)::jsonb").
		From(q.collection.Name)

	if query.Range != nil {
		statBuilder = statBuilder.
			Limit(uint64(query.Range.TakeValue)).
			Offset(uint64(query.Range.SkipValue))
	}

	if len(query.Sorting) > 0 {
		for key, order := range query.Sorting {
			if order.IsAscending {
				statBuilder = statBuilder.OrderBy(key)
			} else {
				statBuilder = statBuilder.OrderBy(key + " DESC")
			}
		}
	}

	condition, err := filtersToCondition("data", query.Filters...)
	if err != nil {
		return nil, err
	}

	if condition != nil {
		statBuilder = statBuilder.Where(condition)
	}

	stat, vals, err := statBuilder.
		ToSql()
	if err != nil {
		return nil, err
	}

	queryCtx, cancel := q.conn.config.GetContext(ctx)

	dbConn, release, err := q.conn.useDatabase(queryCtx, q.collection.Database.Name)
	if err != nil {
		release()
		cancel()
		return nil, err
	}

	rows, err := dbConn.Query(queryCtx, stat, vals...)
	if err != nil {
		release()
		cancel()
		return nil, err
	}

	channel := make(chan T)

	go func(rows pgx.Rows) {
		defer rows.Close()
		defer release()
		defer cancel()
		defer close(channel)

		for rows.Next() {
			select {
			case <-ctx.Done():
				return
			default:
				var data []byte
				if err := rows.Scan(&data); err != nil {
					return
				}

				dataAsBson := make(bson.Raw, 0, len(data))
				dataAsBson, err := jibby.Unmarshal(data, dataAsBson)
				if err != nil {
					return
				}

				var result T
				if err := bson.Unmarshal(dataAsBson, &result); err != nil {
					panic(err)
				}
				channel <- result
			}
		}
	}(rows)

	return channel, nil
}
