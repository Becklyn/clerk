package clerk

type Index struct {
	Fields   []*Field
	Name     string
	IsUnique bool
}

func NewIndex(name ...string) *Index {
	return &Index{
		Fields: []*Field{},
		Name: func() string {
			if len(name) == 0 {
				return ""
			}
			return name[0]
		}(),
	}
}

func (i *Index) AddField(field ...*Field) *Index {
	i.Fields = append(i.Fields, field...)
	return i
}

func (i *Index) Unique() *Index {
	i.IsUnique = true
	return i
}
