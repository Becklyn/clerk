package clerk

type Operation int

const (
	OperationCreate Operation = iota
	OperationUpdate
	OperationDelete
)

var (
	OperationAny = []Operation{
		OperationCreate,
		OperationUpdate,
		OperationDelete,
	}
)

func (o Operation) String() string {
	switch o {
	case OperationCreate:
		return "create"
	case OperationUpdate:
		return "update"
	case OperationDelete:
		return "delete"
	default:
		return "unknown"
	}
}
