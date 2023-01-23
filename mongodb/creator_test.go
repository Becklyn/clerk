package mongodb_test

import (
	"context"
	"testing"

	"github.com/Becklyn/clerk/v4"
	"github.com/Becklyn/clerk/v4/mongodb"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func Test_Creator_CreatesAnEntity(t *testing.T) {
	connection := NewIntegrationConnection(t)

	database := clerk.NewDatabase("integration")
	collection := clerk.NewCollection(database, uuid.NewV4().String())

	type Message struct {
		Id   string `bson:"_id"`
		Text string `bson:"text"`
	}

	message := Message{
		Id:   uuid.NewV4().String(),
		Text: "Hello World",
	}

	operator := mongodb.NewOperator[*Message](connection, collection)

	err := clerk.NewCreate[*Message](operator).
		With(&message).
		Commit(context.Background())
	assert.NoError(t, err)
}

func Test_Creator_CreatesMultipleEntities(t *testing.T) {
	connection := NewIntegrationConnection(t)

	database := clerk.NewDatabase("integration")
	collection := clerk.NewCollection(database, uuid.NewV4().String())

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

	operator := mongodb.NewOperator[*Message](connection, collection)

	err := clerk.NewCreate[*Message](operator).
		With(&message1).
		With(&message2).
		Commit(context.Background())
	assert.NoError(t, err)
}
