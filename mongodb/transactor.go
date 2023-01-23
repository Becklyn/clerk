package mongodb

import (
	"context"

	"github.com/Becklyn/clerk/v4"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type transactor struct {
	connection *Connection
}

func newTransactor(connection *Connection) *transactor {
	return &transactor{
		connection: connection,
	}
}

func (t *transactor) ExecuteTransaction(ctx context.Context, fn clerk.TransactionFn) error {
	transactionCtx, cancel := t.connection.config.GetContext(ctx)
	defer cancel()

	session, err := t.connection.client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(transactionCtx)

	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Snapshot()
	opts := options.Transaction().
		SetWriteConcern(wc).
		SetReadConcern(rc)

	sessionCallback := func(sessCtx mongo.SessionContext) (any, error) {
		return nil, fn(sessCtx)
	}

	_, err = session.WithTransaction(transactionCtx, sessionCallback, opts)
	return err
}
