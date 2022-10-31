package clerk

import "context"

type IndexCreator interface {
	Create(
		ctx context.Context,
		collection *Collection,
		indices ...*IndexCreate,
	) ([]string, error)
}

type IndexDeleter interface {
	Delete(
		ctx context.Context,
		collection *Collection,
		index *IndexDelete,
	) error
}
