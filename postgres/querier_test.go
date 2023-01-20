package postgres_test

import (
	"context"
	"github.com/Becklyn/clerk/v3"
	"github.com/Becklyn/clerk/v3/postgres"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Querier_FindsASingleEntity(t *testing.T) {
	conn := NewIntegrationConnection(t)

	database := clerk.NewDatabase("integration")
	collection := clerk.NewCollection(database, "finds_a_single_entity")

	type Message struct {
		Id   string `bson:"_id"`
		Text string `bson:"text"`
	}

	message := Message{
		Id:   "1",
		Text: "Hello World",
	}

	operator := postgres.NewOperator[*Message](conn, collection)

	result, err := clerk.NewQuery[*Message](operator).
		Where(clerk.NewEquals("_id", message.Id)).
		Single(context.Background())
	assert.NoError(t, err)

	assert.Equal(t, &message, result)
}
