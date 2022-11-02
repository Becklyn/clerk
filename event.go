package clerk

type Event[T any] struct {
	Operation Operation
	Data      T
}
