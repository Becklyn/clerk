package postgres_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Becklyn/clerk/v3"
	"github.com/Becklyn/clerk/v3/postgres"
	"github.com/stretchr/testify/assert"
)

func TestTransactor_Rollback(t *testing.T) {
	conn := postgres.NewIntegrationConnection(t)
	databaseOperator := postgres.NewDatabaseOperator(conn)
	type Message struct {
		Id   string `bson:"_id"`
		Text string `bson:"text"`
	}

	database := clerk.NewDatabase("test_database")
	collection := clerk.NewCollection(database, "test_collection")
	operator := postgres.NewOperator[*Message](conn, collection)

	called := false

	err := clerk.NewTransaction(databaseOperator).Run(context.Background(), func(ctx context.Context) error {
		err := operator.ExecuteCreate(ctx, &clerk.Create[*Message]{
			Data: []*Message{
				{
					Id:   "id",
					Text: "text",
				},
			},
		})
		assert.NoError(t, err)

		called = true
		return errors.New("err")
	})
	assert.Error(t, err)
	assert.True(t, called)

	result, err := clerk.NewQuery[*Message](operator).
		Where(clerk.NewEquals("_id", "id")).
		Single(context.Background())
	assert.NoError(t, err)

	assert.Equal(t, (*Message)(nil), result)
}

func TestTransactor_Commit(t *testing.T) {
	conn := postgres.NewIntegrationConnection(t)
	databaseOperator := postgres.NewDatabaseOperator(conn)
	type Message struct {
		Id   string `bson:"_id"`
		Text string `bson:"text"`
	}

	database := clerk.NewDatabase("test_database")
	collection := clerk.NewCollection(database, "test_collection")
	operator := postgres.NewOperator[*Message](conn, collection)

	called := false

	err := clerk.NewTransaction(databaseOperator).Run(context.Background(), func(ctx context.Context) error {
		err := operator.ExecuteCreate(ctx, &clerk.Create[*Message]{
			Data: []*Message{
				{
					Id:   "other_id",
					Text: "text",
				},
			},
		})
		assert.NoError(t, err)

		called = true
		return nil
	})
	assert.NoError(t, err)
	assert.True(t, called)

	result, err := clerk.NewQuery[*Message](operator).
		Where(clerk.NewEquals("_id", "other_id")).
		Single(context.Background())
	assert.NoError(t, err)

	assert.Equal(t, "other_id", result.Id)
}
