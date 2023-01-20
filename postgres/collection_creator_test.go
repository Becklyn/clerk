package postgres_test

import (
	"testing"
)

func Test_CollectionCreator_CreatesNonExisitingCollection(t *testing.T) {
	// conn := NewIntegrationConnection(t)

	// database := clerk.NewDatabase("test_database")
	// collectionOperator := postgres.NewCollectionOperator(conn, database)

	// collection := clerk.NewCollection(database, "existing_test_database")

	// err := clerk.NewCreate[*clerk.Collection](collectionOperator).
	// 	With(database).
	// 	Commit(context.Background())
	// assert.NoError(t, err)
}

func Test_CollectionCreator_DoesNotCreateExistingCollection(t *testing.T) {
	// conn := NewIntegrationConnection(t)

	// databaseOperator := postgres.NewDatabaseOperator(conn)

	// database := clerk.NewDatabase("existing_test_database")

	// err := clerk.NewCreate[*clerk.Database](databaseOperator).
	// 	With(database).
	// 	Commit(context.Background())
	// assert.NoError(t, err)

	// err = clerk.NewCreate[*clerk.Database](databaseOperator).
	// 	With(database).
	// 	Commit(context.Background())
	// assert.Error(t, err)
}
