package clerk

import "context"

type IndexCreate struct {
	Fields []*IndexField
	Name   string
	Unique bool
}

func NewIndexCreate() *IndexCreate {
	return &IndexCreate{}
}

func (i *IndexCreate) AddField(field *IndexField) *IndexCreate {
	i.Fields = append(i.Fields, field)
	return i
}

func (i *IndexCreate) WithName(name string) *IndexCreate {
	i.Name = name
	return i
}

func (i *IndexCreate) MarkUnique() *IndexCreate {
	i.Unique = true
	return i
}

func (i *IndexCreate) Execute(ctx context.Context, creator IndexCreator) (string, error) {
	return creator.Create(ctx, i)
}
