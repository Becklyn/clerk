package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/Becklyn/clerk/v4"
)

type collectionDeleter struct {
	collectionBase
}

func newCollectionDeleter(conn *Connection, database *clerk.Database) *collectionDeleter {
	return &collectionDeleter{
		*newCollectionBase(conn, database),
	}
}

func (d *collectionDeleter) ExecuteDelete(
	ctx context.Context,
	delete *clerk.Delete[*clerk.Collection],
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

	dbConn, release, err := d.conn.createOrUseDatabase(deleteCtx, d.database.Name)
	defer release()
	if err != nil {
		return 0, err
	}

	for _, name := range names {
		stmt := fmt.Sprintf("DROP TABLE IF EXISTS %s", name)
		if _, err := dbConn.Exec(deleteCtx, stmt); err != nil {
			return 0, err
		}
	}

	return len(names), nil
}
