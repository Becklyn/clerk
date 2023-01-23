package postgres_test

import (
	"context"
	"testing"

	"github.com/Becklyn/clerk/v4"
	"github.com/Becklyn/clerk/v4/postgres"
	"github.com/stretchr/testify/assert"
)

func Test_DatabaseDeleter_DeletesExistingDatabase(t *testing.T) {
	conn := postgres.NewIntegrationConnection(t)

	databaseOperator := postgres.NewDatabaseOperator(conn)

	database := clerk.NewDatabase("existing_test_database")

	err := clerk.NewCreate[*clerk.Database](databaseOperator).
		With(database).
		Commit(context.Background())
	assert.NoError(t, err)

	count, err := clerk.NewDelete[*clerk.Database](databaseOperator).
		Where(clerk.NewEquals("name", database.Name)).
		Commit(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestDatabaseDeleter_DoesNotDeleteNonExistingDatabase(t *testing.T) {
	conn := postgres.NewIntegrationConnection(t)

	databaseOperator := postgres.NewDatabaseOperator(conn)

	count, err := clerk.NewDelete[*clerk.Database](databaseOperator).
		Where(clerk.NewEquals("name", "non_existing_database")).
		Commit(context.Background())
	assert.Error(t, err)
	assert.Equal(t, 0, count)
}
