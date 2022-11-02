package clerk

import "context"

type Create[T any] struct {
	creator Creator[T]
	Data    []T
}

func NewCreate[T any](creator Creator[T]) *Create[T] {
	return &Create[T]{
		creator: creator,
		Data:    []T{},
	}
}

func (c *Create[T]) With(data ...T) *Create[T] {
	c.Data = append(c.Data, data...)
	return c
}

func (c *Create[T]) Commit(ctx context.Context) error {
	return c.creator.ExecuteCreate(ctx, c)
}
