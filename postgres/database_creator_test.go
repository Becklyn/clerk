package postgres_test

import (
	"context"
	"testing"

	"github.com/Becklyn/clerk/v4"
	"github.com/Becklyn/clerk/v4/postgres"
	"github.com/stretchr/testify/assert"
)

func Test_DatabaseCreator_CreatesNonExisitingDatabase(t *testing.T) {
	conn := postgres.NewIntegrationConnection(t)

	databaseOperator := postgres.NewDatabaseOperator(conn)

	database := clerk.NewDatabase("test_database")

	defer func() {
		_, err := clerk.NewDelete[*clerk.Database](databaseOperator).
			Where(clerk.NewEquals("name", database.Name)).
			Commit(context.Background())
		assert.NoError(t, err)
	}()

	err := clerk.NewCreate[*clerk.Database](databaseOperator).
		With(database).
		Commit(context.Background())
	assert.NoError(t, err)
}

func Test_DatabaseCreator_DoesNotCreateExistingDatabase(t *testing.T) {
	conn := postgres.NewIntegrationConnection(t)

	databaseOperator := postgres.NewDatabaseOperator(conn)

	database := clerk.NewDatabase("existing_test_database")

	defer func() {
		_, err := clerk.NewDelete[*clerk.Database](databaseOperator).
			Where(clerk.NewEquals("name", database.Name)).
			Commit(context.Background())
		assert.NoError(t, err)
	}()

	err := clerk.NewCreate[*clerk.Database](databaseOperator).
		With(database).
		Commit(context.Background())
	assert.NoError(t, err)

	err = clerk.NewCreate[*clerk.Database](databaseOperator).
		With(database).
		Commit(context.Background())
	assert.Error(t, err)
}
