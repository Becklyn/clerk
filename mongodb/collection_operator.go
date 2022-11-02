package mongodb

import "github.com/Becklyn/clerk/v2"

type collectionOperator struct {
	collectionQuerier
	collectionDeleter
	collectionUpdater
}

func NewCollectionOperator(connection *Connection, database *clerk.Database) *collectionOperator {
	return &collectionOperator{
		collectionQuerier: *newCollectionQuerier(connection, database),
		collectionDeleter: *newCollectionDeleter(connection, database),
		collectionUpdater: *newCollectionUpdater(connection, database),
	}
}
