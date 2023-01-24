package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/Becklyn/clerk/v4"
)

type indexDeleter struct {
	conn       *Connection
	collection *clerk.Collection
	transactor *transactor
}

func newIndexDeleter(conn *Connection, collection *clerk.Collection) *indexDeleter {
	return &indexDeleter{
		conn:       conn,
		collection: collection,
		transactor: newTransactor(conn),
	}
}

func (d *indexDeleter) ExecuteDelete(
	ctx context.Context,
	delete *clerk.Delete[*clerk.Index],
) (int, error) {
	names := []string{}
	for _, filter := range delete.Filters {
		switch filter.(type) {
		case *clerk.Equals:
			if strings.ToLower(filter.Key()) == "name" {
				names = append(names, filter.Value().(string))
			}
		}
	}

	deleteCtx, cancel := d.conn.config.GetContext(ctx)
	defer cancel()

	if err := d.transactor.ExecuteTransaction(deleteCtx, func(ctx context.Context) error {
		dbConn, release, err := d.conn.createOrUseDatabase(ctx, d.collection.Database.Name)
		defer release()
		if err != nil {
			return err
		}

		if len(names) == 0 {
			rows, err := dbConn.Query(
				ctx,
				"SELECT indexname FROM pg_indexes WHERE tablename = $1 AND indexname != $2",
				d.collection.Name,
				fmt.Sprintf("%s_pk", d.collection.Name),
			)
			if err != nil {
				return err
			}

			var indexNames []string

			for rows.Next() {
				var indexName string
				if err := rows.Scan(&indexName); err != nil {
					rows.Close()
					return err
				}

				indexNames = append(indexNames, indexName)
			}
			rows.Close()

			for _, indexName := range indexNames {
				stmt := fmt.Sprintf("DROP INDEX IF EXISTS %s", indexName)
				if _, err := dbConn.Exec(ctx, stmt); err != nil {
					return err
				}
			}
			return nil
		}

		for _, name := range names {
			stmt := fmt.Sprintf("DROP INDEX IF EXISTS %s", fmt.Sprintf("%s_%s", d.collection.Name, name))
			if _, err := dbConn.Exec(ctx, stmt); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return 0, err
	}

	return len(names), nil
}
