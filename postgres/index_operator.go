package postgres

import "github.com/Becklyn/clerk/v4"

type IndexOperator struct {
	indexQuerier
	indexCreator
	indexDeleter
}

func NewIndexOperator(connection *Connection, collection *clerk.Collection) *IndexOperator {
	return &IndexOperator{
		indexQuerier: *newIndexQuerier(connection, collection),
		indexCreator: *newIndexCreator(connection, collection),
		indexDeleter: *newIndexDeleter(connection, collection),
	}
}
