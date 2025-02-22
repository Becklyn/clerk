package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/Becklyn/clerk/v4"
)

type collectionUpdater struct {
	collectionBase
	transactor *transactor
}

func newCollectionUpdater(conn *Connection, database *clerk.Database) *collectionUpdater {
	return &collectionUpdater{
		collectionBase: *newCollectionBase(conn, database),
		transactor:     newTransactor(conn),
	}
}

func (u *collectionUpdater) ExecuteUpdate(
	ctx context.Context,
	update *clerk.Update[*clerk.Collection],
) error {
	names := []string{}
	for _, filter := range update.Filters {
		switch filter.(type) {
		case *clerk.Equals:
			if strings.ToLower(filter.Key()) == "name" {
				names = append(names, filter.Value().(string))
			}
		}
	}

	updateCtx, cancel := u.conn.config.GetContext(ctx)
	defer cancel()

	return u.transactor.ExecuteTransaction(updateCtx, func(ctx context.Context) error {
		dbConn, release, err := u.conn.createOrUseDatabase(ctx, u.database.Name)
		defer release()
		if err != nil {
			return err
		}

		for _, name := range names {
			rows, err := dbConn.Query(ctx, "SELECT indexname FROM pg_indexes WHERE tablename = $1", name)
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
				suffix := strings.TrimPrefix(indexName, name)
				indexStmt := fmt.Sprintf("ALTER INDEX IF EXISTS %s%s RENAME TO %s%s", name, suffix, update.Data.Name, suffix)
				if _, err := dbConn.Exec(ctx, indexStmt); err != nil {
					return err
				}
			}

			stmt := fmt.Sprintf("ALTER TABLE IF EXISTS %s RENAME TO %s;", name, update.Data.Name)
			if _, err := dbConn.Exec(ctx, stmt); err != nil {
				return err
			}
		}

		return nil
	})
}
