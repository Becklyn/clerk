package postgres

import (
	"context"
	"github.com/Becklyn/clerk/v3"
	"strings"
)

type databaseDeleter struct {
	conn *Connection
}

func newDatabaseDeleter(conn *Connection) *databaseDeleter {
	return &databaseDeleter{
		conn: conn,
	}
}

func (d *databaseDeleter) ExecuteDelete(
	ctx context.Context,
	delete *clerk.Delete[*clerk.Database],
) (int, error) {
	var names []string
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

	for _, name := range names {
		stat := "DROP DATABASE " + name
		if _, err := d.conn.client.Exec(deleteCtx, stat); err != nil {
			return 0, err
		}
	}
	return len(names), nil
}
