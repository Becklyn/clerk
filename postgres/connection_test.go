package postgres

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func isRunningInContainer() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}
	return false
}

func NewIntegrationConnection(t *testing.T) *Connection {
	hostname := "localhost"
	if isRunningInContainer() {
		hostname = "host.docker.internal"
	}

	host := Host(fmt.Sprintf("postgres://postgres:change-me@%s:5432", hostname))

	connection, err := NewConnection(
		context.Background(),
		Config{
			Host:    host,
			Timeout: time.Hour,
		},
	)
	assert.NoError(t, err)
	return connection
}

func TestCanConnectToIntegration(t *testing.T) {
	connection := NewIntegrationConnection(t)
	defer connection.Close()
}
