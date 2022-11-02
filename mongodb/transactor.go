package mongodb

import (
	"context"

	"github.com/Becklyn/clerk/v2"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type transactor struct {
	client *mongo.Client
}

func newTransactor(connection *Connection) *transactor {
	return &transactor{
		client: connection.client,
	}
}

func (t *transactor) ExecuteTransaction(ctx context.Context, fn clerk.TransactionFn) error {
	session, err := t.client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Snapshot()
	opts := options.Transaction().
		SetWriteConcern(wc).
		SetReadConcern(rc)

	sessionCallback := func(sessCtx mongo.SessionContext) (any, error) {
		return nil, fn(sessCtx)
	}

	_, err = session.WithTransaction(ctx, sessionCallback, opts)
	return err
}
