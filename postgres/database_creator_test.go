package postgres_test

import (
	"context"
	"github.com/Becklyn/clerk/v3"
	"github.com/Becklyn/clerk/v3/postgres"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_DatabaseCreator_CreatesNonExisitingDatabase(t *testing.T) {
	conn := NewIntegrationConnection(t)

	databaseOperator := postgres.NewDatabaseOperator(conn)

	database := clerk.NewDatabase("test_database")

	err := clerk.NewCreate[*clerk.Database](databaseOperator).
		With(database).
		Commit(context.Background())
	assert.NoError(t, err)
}

func Test_DatabaseCreator_DoesNotCreateExistingDatabase(t *testing.T) {
	conn := NewIntegrationConnection(t)

	databaseOperator := postgres.NewDatabaseOperator(conn)

	database := clerk.NewDatabase("existing_test_database")

	err := clerk.NewCreate[*clerk.Database](databaseOperator).
		With(database).
		Commit(context.Background())
	assert.NoError(t, err)

	err = clerk.NewCreate[*clerk.Database](databaseOperator).
		With(database).
		Commit(context.Background())
	assert.Error(t, err)
}
