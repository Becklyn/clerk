package postgres

import (
	"context"
	"strings"

	"github.com/Becklyn/clerk/v4"
	sq "github.com/Masterminds/squirrel"
)

type collectionQuerier struct {
	collectionBase
}

func newCollectionQuerier(conn *Connection, database *clerk.Database) *collectionQuerier {
	return &collectionQuerier{
		*newCollectionBase(conn, database),
	}
}

func (q *collectionQuerier) ExecuteQuery(
	ctx context.Context,
	query *clerk.Query[*clerk.Collection],
) (<-chan *clerk.Collection, error) {
	condition, err := filtersToCondition("", query.Filters...)
	if err != nil {
		return nil, err
	}

	filter := sq.And{
		sq.NotEq{
			"schemaname": "pg_catalog",
		},
		sq.NotEq{
			"schemaname": "information_schema",
		},
	}

	if condition != nil {
		filter = append(filter, condition)
	}

	stat, vals, err := statementBuilder().
		Select("name").
		From("pg_catalog.pg_tables").
		Where(filter).
		ToSql()
	if err != nil {
		return nil, err
	}
	stat = strings.ReplaceAll(stat, " name ", " tablename ")

	queryCtx, cancel := q.conn.config.GetContext(ctx)
	defer cancel()

	dbConn, release, err := q.conn.createOrUseDatabase(queryCtx, q.database.Name)
	defer release()
	if err != nil {
		return nil, err
	}

	rows, err := dbConn.Query(queryCtx, stat, vals...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var collections []*clerk.Collection

	for rows.Next() {
		var name string

		err := rows.Scan(&name)
		if err != nil {
			return nil, err
		}

		collections = append(collections, clerk.NewCollection(q.database, name))
	}

	channel := make(chan *clerk.Collection)

	go func() {
		defer close(channel)

		for _, collection := range collections {
			channel <- collection
		}
	}()

	return channel, nil
}
