package postgres_test

import (
	"context"
	"testing"

	"github.com/Becklyn/clerk/v4"
	"github.com/Becklyn/clerk/v4/postgres"
	"github.com/stretchr/testify/assert"
)

func Test_IndexDeleter_DeletesExistingIndex(t *testing.T) {
	conn := postgres.NewIntegrationConnection(t)

	database := clerk.NewDatabase("test_database")
	collection := clerk.NewCollection(database, "test_collection")

	indexOperator := postgres.NewIndexOperator(conn, collection)

	err := clerk.NewCreate[*clerk.Index](indexOperator).
		With(clerk.NewIndex("test_index").AddField(clerk.NewField("field"))).
		Commit(context.Background())
	assert.NoError(t, err)

	count, err := clerk.NewDelete[*clerk.Index](indexOperator).
		Where(clerk.NewEquals("name", "test_index")).
		Commit(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}

func Test_IndexDeleter_DeletesAllIndices(t *testing.T) {
	conn := postgres.NewIntegrationConnection(t)

	database := clerk.NewDatabase("test_database")
	collection := clerk.NewCollection(database, "test_collection")

	indexOperator := postgres.NewIndexOperator(conn, collection)

	count, err := clerk.NewDelete[*clerk.Index](indexOperator).
		Commit(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 0, count)
}
