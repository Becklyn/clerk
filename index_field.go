package clerk

type IndexField struct {
	Key  string
	Type any
}

func NewIndexField(key string) *IndexField {
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
