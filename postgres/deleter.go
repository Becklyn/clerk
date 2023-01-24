package postgres

import (
	"context"

	"github.com/Becklyn/clerk/v4"
)

type deleter[T any] struct {
	conn       *Connection
	collection *clerk.Collection
	transactor *transactor
}

func newDeleter[T any](conn *Connection, collection *clerk.Collection) *deleter[T] {
	return &deleter[T]{
		conn:       conn,
		collection: collection,
		transactor: newTransactor(conn),
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

	var rowsAffected int

	err = d.transactor.ExecuteTransaction(queryCtx, func(ctx context.Context) error {
		dbConn, release, err := d.conn.createOrUseDatabase(ctx, d.collection.Database.Name)
		defer release()
		if err != nil {
			return err
		}

		cmd, err := dbConn.Exec(ctx, stat, vals...)
		if err != nil {
			return err
		}

		rowsAffected = int(cmd.RowsAffected())
		return nil
	})
	return rowsAffected, err
}
