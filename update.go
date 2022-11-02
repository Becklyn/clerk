package clerk

import "context"

type Update[T any] struct {
	updater      Updater[T]
	Filters      []Filter
	ShouldUpsert bool
	Data         T
}

func NewUpdate[T any](updater Updater[T]) *Update[T] {
	var empty T
	return &Update[T]{
		updater:      updater,
		Filters:      []Filter{},
		ShouldUpsert: false,
		Data:         empty,
	}
}

func (u *Update[T]) Where(filter Filter) *Update[T] {
	u.Filters = append(u.Filters, filter)
	return u
}

func (u *Update[T]) Upsert() *Update[T] {
	u.ShouldUpsert = true
	return u
}

func (u *Update[T]) With(data T) *Update[T] {
	u.Data = data
	return u
}

func (u *Update[T]) Commit(ctx context.Context) error {
	return u.updater.ExecuteUpdate(ctx, u)
}
