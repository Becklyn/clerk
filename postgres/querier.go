package postgres

import (
	"context"

	"github.com/Becklyn/clerk/v4"
	"github.com/xdg-go/jibby"
	"go.mongodb.org/mongo-driver/bson"
)

type querier[T any] struct {
	conn              *Connection
	collection        *clerk.Collection
	collectionCreator *collectionCreator
	transactor        *transactor
}

func newQuerier[T any](conn *Connection, collection *clerk.Collection) *querier[T] {
	return &querier[T]{
		conn:              conn,
		collection:        collection,
		collectionCreator: newCollectionCreator(conn, collection.Database),
		transactor:        newTransactor(conn),
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
			keySelector := jsonKeyToSelector("data", key, nil)

			if order.IsAscending {
				statBuilder = statBuilder.OrderBy(keySelector)
			} else {
				statBuilder = statBuilder.OrderBy(keySelector + " DESC")
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
	defer cancel()

	var elements []T

	if err := q.transactor.executeInTransactionIfAvailable(queryCtx, q.collection.Database, func(ctx context.Context) error {
		dbConn, release, err := q.conn.createOrUseDatabase(ctx, q.collection.Database.Name)
		defer release()
		if err != nil {
			return err
		}

		rows, err := dbConn.Query(ctx, stat, vals...)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var data []byte
			if err := rows.Scan(&data); err != nil {
				return err
			}

			dataAsBson := make(bson.Raw, 0, len(data))
			dataAsBson, err := jibby.Unmarshal(data, dataAsBson)
			if err != nil {
				return err
			}

			var result T
			if err := bson.Unmarshal(dataAsBson, &result); err != nil {
				return err
			}
			elements = append(elements, result)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	channel := make(chan T)

	go func() {
		defer close(channel)

		for _, element := range elements {
			channel <- element
		}
	}()

	return channel, nil
}
