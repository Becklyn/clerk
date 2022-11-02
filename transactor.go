package clerk

import "context"

type Transactor interface {
	ExecuteTransaction(ctx context.Context, fn TransactionFn) error
}
