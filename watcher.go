package clerk

import "context"

type Watcher[T any] interface {
	ExecuteWatch(ctx context.Context, watch *Watch[T]) (<-chan *Event[T], error)
}
