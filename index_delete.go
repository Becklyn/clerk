package clerk

import "context"

type IndexDelete struct {
	Name string
}

func NewIndexDelete(name string) *IndexDelete {
	return &IndexDelete{
		Name: name,
	}
}

func (i *IndexDelete) Execute(ctx context.Context, deleter IndexDeleter) error {
	return deleter.Delete(ctx, i)
}
