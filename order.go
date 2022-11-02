package clerk

type Order struct {
	IsAscending bool
}

func NewAscendingOrder() *Order {
	return &Order{
		IsAscending: true,
	}
}

func NewDescendingOrder() *Order {
	return &Order{
		IsAscending: false,
	}
}
