package clerk

import "context"

type IndexCreator interface {
	Create(
		ctx context.Context,
		index *Index,
	) error
}
