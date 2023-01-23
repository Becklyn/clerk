package postgres

import "github.com/Becklyn/clerk/v4"

type CollectionOperator struct {
	collectionCreator
	collectionQuerier
	collectionUpdater
	collectionDeleter
}

func NewCollectionOperator(connection *Connection, database *clerk.Database) *CollectionOperator {
	return &CollectionOperator{
		collectionCreator: *newCollectionCreator(connection, database),
		collectionQuerier: *newCollectionQuerier(connection, database),
		collectionUpdater: *newCollectionUpdater(connection, database),
		collectionDeleter: *newCollectionDeleter(connection, database),
	}
}
