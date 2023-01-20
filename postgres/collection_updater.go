package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/Becklyn/clerk/v3"
)

type collectionUpdater struct {
	collectionBase
}

func newCollectionUpdater(conn *Connection, database *clerk.Database) *collectionUpdater {
	return &collectionUpdater{
		*newCollectionBase(conn, database),
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

	db, release, err := u.tryUseDB(updateCtx)
	defer release()
	if err != nil {
		return err
	}

	for _, name := range names {
		stmt := fmt.Sprintf("ALTER TABLE IF EXISTS %s RENAME TO %s", name, update.Data.Name)
		if _, err := db.client.Exec(updateCtx, stmt); err != nil {
			return err
		}
	}

	return nil
}
