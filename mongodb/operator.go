package mongodb

import "github.com/Becklyn/clerk/v4"

type Operator[T any] struct {
	querier[T]
	creator[T]
	deleter[T]
	updater[T]
	watcher[T]
}

func NewOperator[T any](connection *Connection, collection *clerk.Collection) *Operator[T] {
	return &Operator[T]{
		querier: *newQuerier[T](connection, collection),
		creator: *newCreator[T](connection, collection),
		deleter: *newDeleter[T](connection, collection),
		updater: *newUpdater[T](connection, collection),
		watcher: *newWatcher[T](connection, collection),
	}
}
