package clerk

type Field struct {
	Key  string
	Type FieldType
}

func NewField(key string) *Field {
	return &Field{
		Key:  key,
		Type: FieldTypeAscending,
	}
}

func (f *Field) OfTypeSort(order *Order) *Field {
	if order.IsAscending {
		f.Type = FieldTypeAscending
	} else {
		f.Type = FieldTypeDescending
	}
	return f
}

func (f *Field) OfTypeText() *Field {
	f.Type = FieldTypeString
	return f
}
