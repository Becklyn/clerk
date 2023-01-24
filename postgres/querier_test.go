package postgres_test

import (
	"context"
	"testing"

	"github.com/Becklyn/clerk/v4"
	"github.com/Becklyn/clerk/v4/postgres"
	"github.com/stretchr/testify/assert"
)

func Test_Querier_FindsASingleEntity(t *testing.T) {
	conn := postgres.NewIntegrationConnection(t)

	database := clerk.NewDatabase("test_database")
	collection := clerk.NewCollection(database, "test_collection")

	type Message struct {
		Id   string `bson:"_id"`
		Text string `bson:"text"`
	}

	message := Message{
		Id:   "2",
		Text: "Hello World",
	}

	operator := postgres.NewOperator[*Message](conn, collection)

	result, err := clerk.NewQuery[*Message](operator).
		Where(clerk.NewEquals("_id", message.Id)).
		Single(context.Background())
	assert.NoError(t, err)

	assert.Equal(t, &message, result)
}
