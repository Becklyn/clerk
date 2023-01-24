package postgres

import (
	"context"
	"strings"

	"github.com/Becklyn/clerk/v4"
)

type databaseQuerier struct {
	conn *Connection
}

func newDatabaseQuerier(conn *Connection) *databaseQuerier {
	return &databaseQuerier{
		conn: conn,
	}
}

func (q *databaseQuerier) ExecuteQuery(
	ctx context.Context,
	query *clerk.Query[*clerk.Database],
) (<-chan *clerk.Database, error) {
	condition, err := filtersToCondition("", query.Filters...)
	if err != nil {
		return nil, err
	}

	stat, vals, err := statementBuilder().
		Select("name").
		From("pg_database").
		Where(condition).
		ToSql()
	if err != nil {
		return nil, err
	}
	stat = strings.ReplaceAll(stat, "name", "datname")

	queryCtx, cancel := q.conn.config.GetContext(ctx)
	defer cancel()

	rows, err := q.conn.pool.Query(queryCtx, stat, vals...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var databases []*clerk.Database

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}

		databases = append(databases, clerk.NewDatabase(name))
	}

	channel := make(chan *clerk.Database)

	go func() {
		defer close(channel)

		for _, database := range databases {
			channel <- database
		}
	}()

	return channel, nil
}
