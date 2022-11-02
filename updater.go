package clerk

import "context"

type Updater[T any] interface {
	ExecuteUpdate(ctx context.Context, update *Update[T]) error
}
