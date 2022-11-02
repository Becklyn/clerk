package clerk

import "context"

type Creator[T any] interface {
	ExecuteCreate(ctx context.Context, create *Create[T]) error
}
