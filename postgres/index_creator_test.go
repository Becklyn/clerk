package postgres_test

import (
	"context"
	"sync"
	"testing"

	"github.com/Becklyn/clerk/v4"
	"github.com/Becklyn/clerk/v4/postgres"
	"github.com/stretchr/testify/assert"
)

func Test_IndexCreator_CreatesNonExisitingIndex(t *testing.T) {
	conn := postgres.NewIntegrationConnection(t)

	database := clerk.NewDatabase("test_database")
	collection := clerk.NewCollection(database, "test_collection")

	indexOperator := postgres.NewIndexOperator(conn, collection)

	index := clerk.NewIndex("test").
		AddField(
			clerk.NewField("test").
				OfTypeSort(clerk.NewAscendingOrder()),
			clerk.NewField("test2").
				OfTypeSort(clerk.NewDescendingOrder()),
		).
		Unique()

	err := clerk.NewCreate[*clerk.Index](indexOperator).
		With(index).
		Commit(context.Background())
	assert.NoError(t, err)
}

func Test_IndexCreator_CreateTwiceInTransaction(t *testing.T) {
	conn := postgres.NewIntegrationConnection(t)

	database := clerk.NewDatabase("test_database")
	databaseOperator := postgres.NewDatabaseOperator(conn)
	collection := clerk.NewCollection(database, "test_collection")

	indexOperator := postgres.NewIndexOperator(conn, collection)

	err := clerk.NewTransaction(databaseOperator).Run(context.Background(), func(ctx context.Context) error {
		index := clerk.NewIndex("test").
			AddField(clerk.NewField("test"))

		if err := clerk.NewCreate[*clerk.Index](indexOperator).
			With(index).
			Commit(context.Background()); err != nil {
			return err
		}
		index2 := clerk.NewIndex("test").
			AddField(clerk.NewField("test"))

		return clerk.NewCreate[*clerk.Index](indexOperator).
			With(index2).
			Commit(context.Background())
	})
	assert.NoError(t, err)
}

func Test_IndexCreator_MassiveParallel(t *testing.T) {
	conn := postgres.NewIntegrationConnection(t)

	database := clerk.NewDatabase("test_database")
	collection := clerk.NewCollection(database, "test_collection")

	indexOperator := postgres.NewIndexOperator(conn, collection)

	iterations := 1000
	wg := sync.WaitGroup{}
	wg.Add(iterations)

	for i := 0; i < iterations; i++ {
		go func() {
			index2 := clerk.NewIndex("test").
				AddField(clerk.NewField("test"))

			err := clerk.NewCreate[*clerk.Index](indexOperator).
				With(index2).
				Commit(context.Background())
			assert.NoError(t, err)
			defer wg.Done()
		}()
	}

	wg.Wait()
}

func Test_IndexCreator_CreatesTextIndex_ResultsInError(t *testing.T) {
	conn := postgres.NewIntegrationConnection(t)

	database := clerk.NewDatabase("test_database")
	collection := clerk.NewCollection(database, "test_collection")

	indexOperator := postgres.NewIndexOperator(conn, collection)

	index := clerk.NewIndex("text_test").
		AddField(
			clerk.NewField("test").
				OfTypeText(),
		).
		Unique()

	err := clerk.NewCreate[*clerk.Index](indexOperator).
		With(index).
		Commit(context.Background())
	assert.Error(t, err)
}
