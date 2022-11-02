package clerk

import "context"

type Deleter[T any] interface {
	ExecuteDelete(ctx context.Context, delete *Delete[T]) (int, error)
}
