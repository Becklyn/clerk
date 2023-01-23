package postgres_test

import (
	"context"
	"testing"

	"github.com/Becklyn/clerk/v4"
	"github.com/Becklyn/clerk/v4/postgres"
	"github.com/stretchr/testify/assert"
)

func Test_IndexQuerier_FindsAllIndices(t *testing.T) {
	conn := postgres.NewIntegrationConnection(t)

	database := clerk.NewDatabase("test_database")
	collection := clerk.NewCollection(database, "test_collection")

	indexOperator := postgres.NewIndexOperator(conn, collection)

	indices, err := clerk.NewQuery[*clerk.Index](indexOperator).
		All(context.Background())
	assert.NoError(t, err)
	assert.NotEmpty(t, indices)

}

func Test_DatabaseQuerier_FindsIndexNamedPostgres(t *testing.T) {
	conn := postgres.NewIntegrationConnection(t)

	database := clerk.NewDatabase("test_database")
	collection := clerk.NewCollection(database, "test_collection")

	indexOperator := postgres.NewIndexOperator(conn, collection)

	index, err := clerk.NewQuery[*clerk.Index](indexOperator).
		Where(clerk.NewEquals("name", "test")).
		Single(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, index)
	assert.Equal(t, "test", index.Name)
}
