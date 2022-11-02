package mongodb

import "github.com/Becklyn/clerk/v2"

type operator[T any] struct {
	querier[T]
	creator[T]
	deleter[T]
	updater[T]
	watcher[T]
}

func NewOperator[T any](connection *Connection, collection *clerk.Collection) *operator[T] {
	return &operator[T]{
		querier: *newQuerier[T](connection, collection),
		creator: *newCreator[T](connection, collection),
		deleter: *newDeleter[T](connection, collection),
		updater: *newUpdater[T](connection, collection),
		watcher: *newWatcher[T](connection, collection),
	}
}
