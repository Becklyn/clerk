package postgres

import (
	"context"
	"fmt"

	"github.com/Becklyn/clerk/v3"
)

type deleter[T any] struct {
	conn       *Connection
	collection *clerk.Collection
}

func newDeleter[T any](conn *Connection, collection *clerk.Collection) *deleter[T] {
	return &deleter[T]{
		conn:       conn,
		collection: collection,
	}
}

func (d *deleter[T]) ExecuteDelete(
	ctx context.Context,
	delete *clerk.Delete[T],
) (int, error) {
	statBuilder := statementBuilder().
		Delete(d.collection.Name)

	condition, err := filtersToCondition("data", delete.Filters...)
	if err != nil {
		return 0, err
	}

	if condition != nil {
		statBuilder = statBuilder.Where(condition)
	}

	stat, vals, err := statBuilder.
		ToSql()
	if err != nil {
		return 0, err
	}

	queryCtx, cancel := d.conn.config.GetContext(ctx)
	defer cancel()

	dbConn, release, err := d.conn.useDatabase(queryCtx, d.collection.Database.Name)
	defer release()
	if err != nil {
		return 0, err
	}

	fmt.Println(stat, vals)

	cmd, err := dbConn.Exec(queryCtx, stat, vals...)
	if err != nil {
		return 0, err
	}

	return int(cmd.RowsAffected()), nil
}
