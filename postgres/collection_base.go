package postgres

import (
	"github.com/Becklyn/clerk/v4"
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
