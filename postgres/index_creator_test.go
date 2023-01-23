package postgres_test

import (
	"context"
	"testing"

	"github.com/Becklyn/clerk/v3"
	"github.com/Becklyn/clerk/v3/postgres"
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
