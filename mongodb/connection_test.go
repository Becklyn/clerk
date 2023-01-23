package mongodb_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/Becklyn/clerk/v3/mongodb"
	"github.com/stretchr/testify/assert"
)

func isRunningInContainer() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}
	return false
}

func NewIntegrationConnection(t *testing.T) *mongodb.Connection {
	hostname := "localhost"
	if isRunningInContainer() {
		hostname = "host.docker.internal"
	}

	host := mongodb.Host(fmt.Sprintf("mongodb://root:change-me@%s:27017", hostname))

	connection, err := mongodb.NewConnection(
		context.Background(),
		mongodb.DefaultConfig(host),
	)
	assert.NoError(t, err)
	return connection
}

func TestCanConnectToIntegration(t *testing.T) {
	connection := NewIntegrationConnection(t)
	defer connection.Close(func(err error) {
		assert.NoError(t, err)
	})
}
