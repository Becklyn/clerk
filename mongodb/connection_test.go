package mongodb_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/Becklyn/clerk/v2/mongodb"
	"github.com/stretchr/testify/assert"
)

func isRunningInContainer() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}
	return false
}

func NewIntegrationConnection(t *testing.T) *mongodb.Connection {
	host := "localhost"
	if isRunningInContainer() {
		host = "host.docker.internal"
	}

	connection, err := mongodb.NewConnection(
		context.Background(),
		fmt.Sprintf("mongodb://%s:27017", host),
	)
	assert.NoError(t, err)
	return connection
}
