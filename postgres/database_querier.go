package postgres

import (
	"context"
	"github.com/Becklyn/clerk/v3"
	"strings"
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

	rows, err := q.conn.client.Query(queryCtx, stat, vals...)
	if err != nil {
		return nil, err
	}

	channel := make(chan *clerk.Database)

	go func() {
		defer rows.Close()
		defer close(channel)

		for rows.Next() {
			var name string
			err := rows.Scan(&name)
			if err != nil {
				return
			}

			channel <- clerk.NewDatabase(name)
		}
	}()

	return channel, nil
}
