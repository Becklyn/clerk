package clerk

import "context"

type Querier[T any] interface {
	ExecuteQuery(ctx context.Context, query *Query[T]) (<-chan T, error)
}
