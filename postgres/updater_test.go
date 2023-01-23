package postgres_test

import (
	"context"
	"testing"

	"github.com/Becklyn/clerk/v4"
	"github.com/Becklyn/clerk/v4/postgres"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func Test_Updater_UpdatesAnEntity(t *testing.T) {
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

	updateMessage := Message{
		Id:   message.Id,
		Text: "Foo Bar",
	}

	err = clerk.NewUpdate[*Message](operator).
		Where(clerk.NewEquals("_id", message.Id)).
		With(&updateMessage).
		Commit(context.Background())
	assert.NoError(t, err)

	updatedMessage, err := clerk.NewQuery[*Message](operator).
		Where(clerk.NewEquals("_id", updateMessage.Id)).
		Single(context.Background())
	assert.NoError(t, err)

	assert.Equal(t, &updateMessage, updatedMessage)
}

func Test_Updater_UpsertsAnEntity(t *testing.T) {
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

	upsertMessage := Message{
		Id:   message.Id,
		Text: "Foo Bar",
	}

	err := clerk.NewUpdate[*Message](operator).
		Where(clerk.NewEquals("_id", message.Id)).
		With(&upsertMessage).
		Upsert().
		Commit(context.Background())
	assert.NoError(t, err)

	upsertedMessage, err := clerk.NewQuery[*Message](operator).
		Where(clerk.NewEquals("_id", upsertMessage.Id)).
		Single(context.Background())
	assert.NoError(t, err)

	assert.Equal(t, &upsertMessage, upsertedMessage)
}
