package postgres_test

import (
	"context"
	"sync"
	"testing"

	"github.com/Becklyn/clerk/v4"
	"github.com/Becklyn/clerk/v4/postgres"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func Test_Operator_MassiveParallel(t *testing.T) {
	conn := postgres.NewIntegrationConnection(t)

	database := clerk.NewDatabase("test_database")
	databaseOperator := postgres.NewDatabaseOperator(conn)
	collection := clerk.NewCollection(database, "test_collection")

	type Message struct {
		Id   string `bson:"_id"`
		Text string `bson:"text"`
	}

	message := Message{
		Id:   "2",
		Text: "Hello World",
	}

	operator := postgres.NewOperator[*Message](conn, collection)

	iterations := 1000
	wg := sync.WaitGroup{}
	wg.Add(iterations)

	for i := 0; i < iterations; i++ {
		go func() {
			err := clerk.NewTransaction(databaseOperator).Run(context.Background(), func(ctx context.Context) error {
				_, err := clerk.NewQuery[*Message](operator).
					Where(clerk.NewEquals("_id", message.Id)).
					Single(context.Background())
				if err != nil {
					return err
				}

				err = clerk.NewCreate[*Message](operator).
					With(&Message{Id: uuid.NewV4().String()}).
					Commit(context.Background())
				return err
			})
			assert.NoError(t, err)
			defer wg.Done()
		}()
	}

	wg.Wait()
}
