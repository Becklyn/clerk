package postgres

import (
	"context"

	"github.com/Becklyn/clerk/v3"
)

type collectionBase struct {
	conn      *Connection
	database  *clerk.Database
	dbCreator *databaseCreator
}

func newCollectionBase(conn *Connection, database *clerk.Database) *collectionBase {
	return &collectionBase{
		conn:      conn,
		database:  database,
		dbCreator: newDatabaseCreator(conn),
	}
}

func (c *collectionBase) tryUseDB(ctx context.Context) (*DatabaseConnection, func(), error) {
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
