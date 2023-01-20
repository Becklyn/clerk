package postgres

import (
	"context"
	"fmt"

	"github.com/Becklyn/clerk/v3"
)

type collectionCreator struct {
	conn      *Connection
	database  *clerk.Database
	dbCreator *databaseCreator
}

func newCollectionCreator(conn *Connection, database *clerk.Database) *collectionCreator {
	return &collectionCreator{
		conn:      conn,
		database:  database,
		dbCreator: newDatabaseCreator(conn),
	}
}

func (c *collectionCreator) ExecuteCreate(
	ctx context.Context,
	create *clerk.Create[*clerk.Collection],
) error {
	createCtx, cancel := c.conn.config.GetContext(ctx)
	defer cancel()

	db, release, err := c.tryUseDB(ctx)
	defer release()
	if err != nil {
		return err
	}

	for _, data := range create.Data {
		stat := fmt.Sprintf("create table %s (id uuid constraint %s_pk primary key, data jsonb default '{}')", data.Name, data.Name)

		_, err := db.Client().Exec(createCtx, stat)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *collectionCreator) tryUseDB(ctx context.Context) (*DatabaseConnection, func(), error) {
	db, release, err := c.conn.UseDatabase(c.database.Name)
	if err != nil {
		if errCreate := c.dbCreator.ExecuteCreate(ctx, &clerk.Create[*clerk.Database]{
			Data: []*clerk.Database{c.database},
		}); errCreate != nil {
			return nil, nil, err
		}

		return c.conn.UseDatabase(c.database.Name)
	}

	return db, release, err
}
