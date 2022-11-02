package clerk

type FieldType int

const (
	FieldTypeAscending FieldType = iota
	FieldTypeDescending
	FieldTypeString
)

func (t FieldType) String() string {
	switch t {
	case FieldTypeAscending:
		return "ascending"
	case FieldTypeDescending:
		return "descending"
	case FieldTypeString:
		return "string"
	default:
		return "unknown"
	}
}
