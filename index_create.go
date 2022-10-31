package clerk

import "context"

type IndexCreate struct {
	Collection *Collection
	Fields     []*IndexField
	Name       string
	Unique     bool
}

func NewIndexCreate(collection *Collection) *IndexCreate {
	return &IndexCreate{
		Collection: collection,
	}
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
	names, err := creator.Create(ctx, i.Collection, i)
	if err != nil {
		return "", err
	}
	return names[0], nil
}
