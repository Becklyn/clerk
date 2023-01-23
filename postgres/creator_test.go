package postgres_test

import (
	"context"
	"testing"

	"github.com/Becklyn/clerk/v3"
	"github.com/Becklyn/clerk/v3/postgres"
	"github.com/stretchr/testify/assert"
)

func Test_Creator_CreatesData(t *testing.T) {
	conn := postgres.NewIntegrationConnection(t)

	type Message struct {
		Id   string `bson:"_id"`
		Text string `bson:"text"`
	}

	database := clerk.NewDatabase("test_database")
	collection := clerk.NewCollection(database, "test_collection")
	collectionOperator := postgres.NewOperator[*Message](conn, collection)

	message := Message{
		Id:   "2",
		Text: "Hello World",
	}

	err := clerk.NewCreate[*Message](collectionOperator).
		With(&message).
		Commit(context.Background())
	assert.NoError(t, err)
}
