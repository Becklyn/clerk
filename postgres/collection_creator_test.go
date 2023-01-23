package postgres_test

import (
	"context"
	"testing"

	"github.com/Becklyn/clerk/v3"
	"github.com/Becklyn/clerk/v3/postgres"
	"github.com/stretchr/testify/assert"
)

func Test_CollectionCreator_CreatesNonExisitingCollection(t *testing.T) {
	conn := postgres.NewIntegrationConnection(t)

	database := clerk.NewDatabase("test_database")
	collectionOperator := postgres.NewCollectionOperator(conn, database)

	collection := clerk.NewCollection(database, "test_collection")

	err := clerk.NewCreate[*clerk.Collection](collectionOperator).
		With(collection).
		Commit(context.Background())
	assert.NoError(t, err)
}

func Test_CollectionCreator_DoesNotCreateExistingCollection(t *testing.T) {
	conn := postgres.NewIntegrationConnection(t)

	database := clerk.NewDatabase("test_database")
	collectionOperator := postgres.NewCollectionOperator(conn, database)

	collection := clerk.NewCollection(database, "existing_test_collection")

	err := clerk.NewCreate[*clerk.Collection](collectionOperator).
		With(collection).
		Commit(context.Background())
	assert.NoError(t, err)

	err = clerk.NewCreate[*clerk.Collection](collectionOperator).
		With(collection).
		Commit(context.Background())
	assert.NoError(t, err)
}
