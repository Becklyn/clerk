package postgres

import "github.com/Becklyn/clerk/v3"

type CollectionOperator struct {
	collectionCreator
}

func NewCollectionOperator(connection *Connection, database *clerk.Database) *CollectionOperator {
	return &CollectionOperator{
		collectionCreator: *newCollectionCreator(connection, database),
	}
}
