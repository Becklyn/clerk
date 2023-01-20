package postgres

import "github.com/Becklyn/clerk/v3"

type Operator[T any] struct {
	querier[T]
}

func NewOperator[T any](conn *Connection, collection *clerk.Collection) *Operator[T] {
	return &Operator[T]{
		querier: *newQuerier[T](conn, collection),
	}
}
