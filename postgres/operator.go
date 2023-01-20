package postgres

import "github.com/Becklyn/clerk/v3"

type Operator[T any] struct {
	querier[T]
	creator[T]
}

func NewOperator[T any](conn *Connection, collection *clerk.Collection) *Operator[T] {
	return &Operator[T]{
		querier: *newQuerier[T](conn, collection),
		creator: *newCreator[T](conn, collection),
	}
}
