package clerk

import "context"

type Index struct {
	Fields []*IndexField
	Name   string
	Unique bool
}

func NewIndex() *Index {
	return &Index{}
}

func (i *Index) AddField(field *IndexField) *Index {
	i.Fields = append(i.Fields, field)
	return i
}

func (i *Index) WithName(name string) *Index {
	i.Name = name
	return i
}

func (i *Index) MarkUnique() *Index {
	i.Unique = true
	return i
}

func (i *Index) Execute(ctx context.Context, creator IndexCreator) error {
	return creator.Create(ctx, i)
}

type IndexField struct {
	Key  string
	Type any
}

func NewField(key string) *IndexField {
	return &IndexField{
		Key: key,
	}
}

func (f *IndexField) OfTypeSort(asc bool) *IndexField {
	if asc {
		f.Type = 1
	} else {
		f.Type = -1
	}
	return f
}

func (f *IndexField) OfTypeText() *IndexField {
	f.Type = "text"
	return f
}
