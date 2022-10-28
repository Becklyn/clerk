package clerk

import "context"

type IndexCreator interface {
	Create(
		ctx context.Context,
		index *IndexCreate,
	) error
}

type IndexDeleter interface {
	Delete(
		ctx context.Context,
		index *IndexDelete,
	) error
}
