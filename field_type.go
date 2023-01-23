package clerk

type FieldType int

const (
	FieldTypeAscending FieldType = iota
	FieldTypeDescending
	FieldTypeText
)

func (t FieldType) String() string {
	switch t {
	case FieldTypeAscending:
		return "ascending"
	case FieldTypeDescending:
		return "descending"
	case FieldTypeText:
		return "text"
	default:
		return "unknown"
	}
}
