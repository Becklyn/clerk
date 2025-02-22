package postgres_test

import (
	"context"
	"testing"

	"github.com/Becklyn/clerk/v4"
	"github.com/Becklyn/clerk/v4/postgres"
	"github.com/stretchr/testify/assert"
)

func Test_DatabaseQuerier_Count(t *testing.T) {
	conn := postgres.NewIntegrationConnection(t)

	databaseOperator := postgres.NewDatabaseOperator(conn)

	total, err := clerk.NewQuery[*clerk.Database](databaseOperator).
		Count(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
}

func Test_DatabaseQuerier_FindsAllDatabases(t *testing.T) {
	conn := postgres.NewIntegrationConnection(t)

	databaseOperator := postgres.NewDatabaseOperator(conn)

	databases, err := clerk.NewQuery[*clerk.Database](databaseOperator).
		All(context.Background())
	assert.NoError(t, err)
	assert.NotEmpty(t, databases)
}

func Test_DatabaseQuerier_FindsDatabaseNamedPostgres(t *testing.T) {
	conn := postgres.NewIntegrationConnection(t)

	databaseOperator := postgres.NewDatabaseOperator(conn)

	database, err := clerk.NewQuery[*clerk.Database](databaseOperator).
		Where(clerk.NewEquals("name", "postgres")).
		Single(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, database)
	assert.Equal(t, "postgres", database.Name)
}
