package postgres_test

import (
	"context"
	"testing"

	"github.com/Becklyn/clerk/v3"
	"github.com/Becklyn/clerk/v3/postgres"
	"github.com/stretchr/testify/assert"
)

func Test_CollectionDeleter_DeleterCollection(t *testing.T) {
	conn := postgres.NewIntegrationConnection(t)

	database := clerk.NewDatabase("test_database")
	collectionOperator := postgres.NewCollectionOperator(conn, database)

	number, err := clerk.NewDelete[*clerk.Collection](collectionOperator).
		Where(clerk.NewEquals("name", "test_collection")).
		Commit(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 1, number)
}

func Test_CollectionDeleter_DeleterNonExistingCollection(t *testing.T) {
	conn := postgres.NewIntegrationConnection(t)

	database := clerk.NewDatabase("test_database")
	collectionOperator := postgres.NewCollectionOperator(conn, database)

	number, err := clerk.NewDelete[*clerk.Collection](collectionOperator).
		Where(clerk.NewEquals("name", "other")).
		Commit(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 1, number)
}
