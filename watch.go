package clerk

import "context"

type Watch[T any] struct {
	watcher    Watcher[T]
	Operations []Operation
}

func NewWatch[T any](watcher Watcher[T]) *Watch[T] {
	return &Watch[T]{
		watcher:    watcher,
		Operations: []Operation{},
	}
}

func (w *Watch[T]) On(operation ...Operation) *Watch[T] {
	w.Operations = append(w.Operations, operation...)
	return w
}

func (w *Watch[T]) Channel(ctx context.Context) (<-chan *Event[T], error) {
	return w.watcher.ExecuteWatch(ctx, w)
}

type WatchHandler[T any] func(ctx context.Context, event *Event[T])

func (w *Watch[T]) Handle(ctx context.Context, handler WatchHandler[T]) error {
	channel, err := w.Channel(ctx)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-channel:
				handler(ctx, event)
			}
		}
	}()

	return nil
}
