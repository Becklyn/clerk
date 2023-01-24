package postgres

import (
	"context"

	"github.com/Becklyn/clerk/v4"
	"github.com/jackc/pgx/v5"
	"go.mongodb.org/mongo-driver/bson"
)

type creator[T any] struct {
	collectionBase
	collectionCreator *collectionCreator
	collection        *clerk.Collection
	transactor        *transactor
}

func newCreator[T any](conn *Connection, collection *clerk.Collection) *creator[T] {
	return &creator[T]{
		collectionBase:    *newCollectionBase(conn, collection.Database),
		collectionCreator: newCollectionCreator(conn, collection.Database),
		collection:        collection,
		transactor:        newTransactor(conn),
	}
}

func (c *creator[T]) ExecuteCreate(
	ctx context.Context,
	create *clerk.Create[T],
) error {
	data := make([]any, len(create.Data))
	for i, item := range create.Data {
		data[i] = item
	}

	createCtx, cancel := c.conn.config.GetContext(ctx)
	defer cancel()

	return c.transactor.ExecuteTransaction(createCtx, func(ctx context.Context) error {
		dbConn, release, err := c.conn.createOrUseDatabase(ctx, c.database.Name)
		defer release()
		if err != nil {
			return err
		}

		for _, data := range create.Data {
			if err := c.create(ctx, data, dbConn); err != nil {
				return err
			}
		}

		return nil
	})
}

func (c *creator[T]) create(ctx context.Context, data T, dbConn *pgx.Conn) error {
	bytes, err := bson.Marshal(data)
	if err != nil {
		return err
	}

	var dataMap map[string]any

	if err := bson.Unmarshal(bytes, &dataMap); err != nil {
		return err
	}

	stat, vals, err := statementBuilder().
		Insert(c.collection.Name).
		Columns("data").
		Values(dataMap).
		ToSql()
	if err != nil {
		return err
	}

	_, err = dbConn.Exec(ctx, stat, vals...)
	return err
}
