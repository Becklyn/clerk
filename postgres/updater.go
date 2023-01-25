package postgres

import (
	"context"
	"encoding/json"

	"github.com/Becklyn/clerk/v4"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"go.mongodb.org/mongo-driver/bson"
)

type updater[T any] struct {
	conn              *Connection
	collection        *clerk.Collection
	collectionCreator *collectionCreator
	transactor        *transactor
}

func newUpdater[T any](conn *Connection, collection *clerk.Collection) *updater[T] {
	return &updater[T]{
		conn:              conn,
		collection:        collection,
		collectionCreator: newCollectionCreator(conn, collection.Database),
		transactor:        newTransactor(conn),
	}
}

func (u *updater[T]) ExecuteUpdate(ctx context.Context, update *clerk.Update[T]) error {
	bytes, err := bson.Marshal(update.Data)
	if err != nil {
		return err
	}

	var dataMap map[string]any

	if err := bson.Unmarshal(bytes, &dataMap); err != nil {
		return err
	}

	condition, err := filtersToCondition("data", update.Filters...)
	if err != nil {
		return err
	}

	updateCtx, cancel := u.conn.config.GetContext(ctx)
	defer cancel()

	return u.transactor.executeInTransactionIfAvailable(updateCtx, u.collection.Database, func(ctx context.Context) error {
		dbConn, release, err := u.conn.createOrUseDatabase(ctx, u.collection.Database.Name)
		defer release()
		if err != nil {
			return err
		}

		if update.ShouldUpsert {
			return u.upsertData(ctx, dbConn, dataMap, condition)
		}

		return u.updateData(ctx, dbConn, dataMap, condition)
	})
}

func (u *updater[T]) upsertData(ctx context.Context, dbConn *pgx.Conn, dataMap map[string]any, condition squirrel.Sqlizer) error {
	dataBytes, err := json.Marshal(dataMap)
	if err != nil {
		return err
	}

	stat, vals, err := statementBuilder().
		Insert(u.collection.Name).
		Columns("data").
		Values(dataMap).
		Suffix("ON CONFLICT ((data->>'_id')) DO UPDATE SET data = ?", string(dataBytes)).
		ToSql()
	if err != nil {
		return err
	}

	_, err = dbConn.Exec(ctx, stat, vals...)
	return err
}

func (u *updater[T]) updateData(ctx context.Context, dbConn *pgx.Conn, dataMap map[string]any, condition squirrel.Sqlizer) error {
	statBuilder := statementBuilder().
		Update(u.collection.Name).
		Set("data", dataMap)

	if condition != nil {
		statBuilder = statBuilder.Where(condition)
	} else {
		statBuilder = statBuilder.Where("SELECT data->'_id' FROM test_collection LIMIT 1")
	}

	stat, vals, err := statBuilder.
		ToSql()
	if err != nil {
		return err
	}

	_, err = dbConn.Exec(ctx, stat, vals...)
	return err
}
