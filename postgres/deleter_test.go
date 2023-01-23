package postgres_test

import (
	"context"
	"testing"

	"github.com/Becklyn/clerk/v3"
	"github.com/Becklyn/clerk/v3/postgres"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func Test_Deleter_DeletesAnEntity(t *testing.T) {
	connection := postgres.NewIntegrationConnection(t)

	database := clerk.NewDatabase("test_database")
	collection := clerk.NewCollection(database, "test_collection")

	type Message struct {
		Id   string `bson:"_id"`
		Text string `bson:"text"`
	}

	message := Message{
		Id:   uuid.NewV4().String(),
		Text: "Hello World",
	}

	operator := postgres.NewOperator[*Message](connection, collection)

	err := clerk.NewCreate[*Message](operator).
		With(&message).
		Commit(context.Background())
	assert.NoError(t, err)

	count, err := clerk.NewDelete[*Message](operator).
		Where(clerk.NewEquals("_id", message.Id)).
		Commit(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}

func Test_Deleter_WithoutAnyFilters_DoesDeletesAllEntities(t *testing.T) {
	connection := postgres.NewIntegrationConnection(t)

	database := clerk.NewDatabase("test_database")
	collection := clerk.NewCollection(database, "test_collection")

	type Message struct {
		Id   string `bson:"_id"`
		Text string `bson:"text"`
	}

	message1 := Message{
		Id:   uuid.NewV4().String(),
		Text: "Hello World",
	}

	message2 := Message{
		Id:   uuid.NewV4().String(),
		Text: "Foo Bar",
	}

	operator := postgres.NewOperator[*Message](connection, collection)

	err := clerk.NewCreate[*Message](operator).
		With(&message1).
		With(&message2).
		Commit(context.Background())
	assert.NoError(t, err)

	count, err := clerk.NewDelete[*Message](operator).
		Commit(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 2, count)
}
