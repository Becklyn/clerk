package clerk

import "context"

type TransactionFn func(ctx context.Context) error

type Transaction struct {
	transactor Transactor
}

func NewTransaction(transactor Transactor) *Transaction {
	return &Transaction{
		transactor: transactor,
	}
}

func (t *Transaction) Run(ctx context.Context, fn TransactionFn) error {
	return t.transactor.ExecuteTransaction(ctx, fn)
}
