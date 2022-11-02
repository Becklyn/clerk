package clerk

import "context"

type Delete[T any] struct {
	deleter Deleter[T]
	Filters []Filter
}

func NewDelete[T any](deleter Deleter[T]) *Delete[T] {
	return &Delete[T]{
		deleter: deleter,
		Filters: []Filter{},
	}
}

func (d *Delete[T]) Where(filter Filter) *Delete[T] {
	d.Filters = append(d.Filters, filter)
	return d
}

func (d *Delete[T]) Commit(ctx context.Context) (int, error) {
	return d.deleter.ExecuteDelete(ctx, d)
}
