package mongodb

import "github.com/Becklyn/clerk/v2"

type CollectionOperator struct {
	collectionQuerier
	collectionDeleter
	collectionUpdater
}

func NewCollectionOperator(connection *Connection, database *clerk.Database) *CollectionOperator {
	return &CollectionOperator{
		collectionQuerier: *newCollectionQuerier(connection, database),
		collectionDeleter: *newCollectionDeleter(connection, database),
		collectionUpdater: *newCollectionUpdater(connection, database),
	}
}
