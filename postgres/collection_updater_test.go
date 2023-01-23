package postgres_test

import (
	"context"
	"testing"

	"github.com/Becklyn/clerk/v3"
	"github.com/Becklyn/clerk/v3/postgres"
	"github.com/stretchr/testify/assert"
)

func Test_CollectionUpdater_UpdateCollection(t *testing.T) {
	conn := postgres.NewIntegrationConnection(t)

	database := clerk.NewDatabase("test_database")
	collectionOperator := postgres.NewCollectionOperator(conn, database)

	collection := clerk.NewCollection(database, "new_test_collection")

	err := clerk.NewUpdate[*clerk.Collection](collectionOperator).
		With(collection).
		Where(clerk.NewEquals("name", "test_collection2")).
		Commit(context.Background())
	assert.NoError(t, err)
}

func Test_CollectionUpdater_UpdateNonExistingCollection(t *testing.T) {
	conn := postgres.NewIntegrationConnection(t)

	database := clerk.NewDatabase("test_database")
	collectionOperator := postgres.NewCollectionOperator(conn, database)

	collection := clerk.NewCollection(database, "test_collection")

	err := clerk.NewUpdate[*clerk.Collection](collectionOperator).
		With(collection).
		Where(clerk.NewEquals("name", "other")).
		Commit(context.Background())
	assert.NoError(t, err)
}
