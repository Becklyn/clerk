package postgres_test

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/Becklyn/clerk/v4"
	"github.com/Becklyn/clerk/v4/postgres"
	uuid "github.com/satori/go.uuid"
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
	id := uuid.NewV4().String()

	err := clerk.NewTransaction(databaseOperator).Run(context.Background(), func(ctx context.Context) error {
		called = true
		return operator.ExecuteCreate(ctx, &clerk.Create[*Message]{
			Data: []*Message{
				{
					Id:   id,
					Text: "text",
				},
			},
		})
	})
	assert.NoError(t, err)
	assert.True(t, called)

	result, err := clerk.NewQuery[*Message](operator).
		Where(clerk.NewEquals("_id", id)).
		Single(context.Background())
	assert.NoError(t, err)

	assert.Equal(t, id, result.Id)
}

func TestTransactor_Nested(t *testing.T) {
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
	id := uuid.NewV4().String()

	err := clerk.NewTransaction(databaseOperator).Run(context.Background(), func(ctx context.Context) error {
		return clerk.NewTransaction(databaseOperator).Run(ctx, func(ctx context.Context) error {
			called = true
			return operator.ExecuteCreate(ctx, &clerk.Create[*Message]{
				Data: []*Message{
					{
						Id:   id,
						Text: "text",
					},
				},
			})
		})
	})
	assert.NoError(t, err)
	assert.True(t, called)

	result, err := clerk.NewQuery[*Message](operator).
		Where(clerk.NewEquals("_id", id)).
		Single(context.Background())
	assert.NoError(t, err)

	assert.Equal(t, id, result.Id)
}

func TestTransactor_ReUse(t *testing.T) {
	conn := postgres.NewIntegrationConnection(t)
	databaseOperator := postgres.NewDatabaseOperator(conn)
	type Message struct {
		Id   string `bson:"_id"`
		Text string `bson:"text"`
	}

	database := clerk.NewDatabase("test_database")
	collection := clerk.NewCollection(database, "test_collection")
	operator := postgres.NewOperator[*Message](conn, collection)

	id := uuid.NewV4().String()

	wg := sync.WaitGroup{}
	wg.Add(100)

	for i := 0; i < 100; i++ {
		go func() {
			called := false
			err := clerk.NewTransaction(databaseOperator).Run(context.Background(), func(ctx context.Context) error {
				called = true
				return operator.ExecuteUpdate(ctx, &clerk.Update[*Message]{
					Data: &Message{
						Id:   id,
						Text: "text",
					},
					ShouldUpsert: true,
					Filters: []clerk.Filter{
						clerk.NewEquals("_id", id),
					},
				})
			})
			assert.NoError(t, err)
			assert.True(t, called)
			wg.Done()
		}()
	}

	wg.Wait()

	results, err := clerk.NewQuery[*Message](operator).
		Where(clerk.NewEquals("_id", id)).
		All(context.Background())
	assert.NoError(t, err)

	assert.Len(t, results, 1)
}
