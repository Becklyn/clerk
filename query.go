package clerk

import "context"

type Query[T any] struct {
	querier Querier[T]
	Filters []Filter
	Sorting map[string]*Order
	Range   *Range
}

func NewQuery[T any](querier Querier[T]) *Query[T] {
	return &Query[T]{
		querier: querier,
		Filters: []Filter{},
		Sorting: map[string]*Order{},
		Range:   nil,
	}
}

func (q *Query[T]) Where(filter Filter) *Query[T] {
	q.Filters = append(q.Filters, filter)
	return q
}

func (q *Query[T]) Sort(key string, order ...*Order) *Query[T] {
	if len(order) > 0 {
		q.Sorting[key] = order[0]
	} else {
		q.Sorting[key] = NewAscendingOrder()
	}
	return q
}

func (q *Query[T]) In(rng *Range) *Query[T] {
	q.Range = rng
	return q
}

func (q *Query[T]) Channel(ctx context.Context) (<-chan T, error) {
	return q.querier.ExecuteQuery(ctx, q)
}

func (q *Query[T]) All(ctx context.Context) ([]T, error) {
	channel, err := q.Channel(ctx)
	if err != nil {
		return nil, err
	}
	var results []T
	for result := range channel {
		results = append(results, result)
	}
	return results, nil
}

func (q *Query[T]) Single(ctx context.Context) (T, error) {
	channel, err := q.Channel(ctx)
	if err != nil {
		var empty T
		return empty, err
	}
	return <-channel, nil
}
