package postgres_test

import (
	"context"
	"testing"

	"github.com/Becklyn/clerk/v4"
	"github.com/Becklyn/clerk/v4/postgres"
	"github.com/stretchr/testify/assert"
)

func Test_CollectionQuerier_Count(t *testing.T) {
	conn := postgres.NewIntegrationConnection(t)

	database := clerk.NewDatabase("test_database")
	collectionOperator := postgres.NewCollectionOperator(conn, database)

	total, err := clerk.NewQuery[*clerk.Collection](collectionOperator).
		Count(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
}

func Test_CollectionQuerier_FindsAllCollections(t *testing.T) {
	conn := postgres.NewIntegrationConnection(t)

	database := clerk.NewDatabase("test_database")
	collectionOperator := postgres.NewCollectionOperator(conn, database)

	collections, err := clerk.NewQuery[*clerk.Collection](collectionOperator).
		All(context.Background())
	assert.NoError(t, err)
	assert.NotEmpty(t, collections)
}

func Test_CollectionQuerier_FindsCollectionNamedPostgres(t *testing.T) {
	conn := postgres.NewIntegrationConnection(t)

	database := clerk.NewDatabase("test_database")
	collectionOperator := postgres.NewCollectionOperator(conn, database)

	collection, err := clerk.NewQuery[*clerk.Collection](collectionOperator).
		Where(clerk.NewEquals("name", "test_collection")).
		Single(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, collection)
	assert.Equal(t, "test_collection", collection.Name)
}
