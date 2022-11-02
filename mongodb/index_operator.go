package mongodb

import "github.com/Becklyn/clerk/v2"

type indexOperator struct {
	indexQuerier
	indexCreator
	indexDeleter
}

func NewIndexOperator(connection *Connection, collection *clerk.Collection) *indexOperator {
	return &indexOperator{
		indexQuerier: *newIndexQuerier(connection, collection),
		indexCreator: *newIndexCreator(connection, collection),
		indexDeleter: *newIndexDeleter(connection, collection),
	}
}
