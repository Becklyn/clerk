package clerk

type Index struct {
	Collection *Collection
	Fields     []*IndexField
	Name       string
	Unique     bool
}
