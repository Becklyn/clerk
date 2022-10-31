package clerk

import "context"

type IndexDelete struct {
	Collection *Collection
	Name       string
}

func NewIndexDelete(collection *Collection, name string) *IndexDelete {
	return &IndexDelete{
		Collection: collection,
		Name:       name,
	}
}

func (i *IndexDelete) Execute(ctx context.Context, deleter IndexDeleter) error {
	return deleter.Delete(ctx, i.Collection, i)
}
