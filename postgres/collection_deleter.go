package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/Becklyn/clerk/v3"
)

type collectionDeleter struct {
	collectionBase
}

func newCollectionDeleter(conn *Connection, database *clerk.Database) *collectionDeleter {
	return &collectionDeleter{
		*newCollectionBase(conn, database),
	}
}

func (u *collectionDeleter) ExecuteDelete(
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

	deleteCtx, cancel := u.conn.config.GetContext(ctx)
	defer cancel()

	db, release, err := u.tryUseDB(deleteCtx)
	defer release()
	if err != nil {
		return 0, err
	}

	for _, name := range names {
		stmt := fmt.Sprintf("DROP TABLE IF EXISTS %s", name)
		if _, err := db.client.Exec(deleteCtx, stmt); err != nil {
			return 0, err
		}
	}

	return len(names), nil
}
