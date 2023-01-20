package postgres_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/Becklyn/clerk/v3/postgres"
	sq "github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/assert"
)

func isRunningInContainer() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}
	return false
}

func NewIntegrationConnection(t *testing.T) *postgres.Connection {
	hostname := "localhost"
	if isRunningInContainer() {
		hostname = "host.docker.internal"
	}

	host := postgres.Host(fmt.Sprintf("postgres://postgres:change-me@%s:5432", hostname))

	connection, err := postgres.NewConnection(
		context.Background(),
		postgres.DefaultConfig(host),
	)
	assert.NoError(t, err)
	return connection
}

func TestSQl(t *testing.T) {
	// CREATE INDEX ON publishers((info->>'name'));q
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	q, p, err := psql.
		Select("*").
		From("test").
		Where(sq.Expr("data @> ''")).
		ToSql()
	assert.NoError(t, err)

	// a := sq.Or{
	// 	sq.Eq{
	// 		"fieldName": 123,
	// 	},
	// 	sq.Eq{
	// 		"fieldName2": 123,
	// 	},
	// }

	// fmt.Println(a.ToSql())

	// assert.False(t, true)

	conn := NewIntegrationConnection(t)
	defer conn.Close(func(err error) {
		assert.NoError(t, err)
	})

	jsondb, release, err := conn.UseDatabase("jsondb")
	defer release()
	assert.NoError(t, err)

	// "SELECT FROM pg_database WHERE datname = 'test'"
	rows, err := jsondb.Client().Query(
		context.Background(),
		q,
		p...,
	)
	assert.NoError(t, err)
	for rows.Next() {
		var id string
		var data string
		err := rows.Scan(&id, &data)
		assert.NoError(t, err)
		fmt.Println(id, data)
	}
	rows.Close()
	assert.False(t, true)
}

func TestCanConnectToIntegration(t *testing.T) {
	connection := NewIntegrationConnection(t)
	defer connection.Close(func(err error) {
		assert.NoError(t, err)
	})
}

func TestSandbox(t *testing.T) {
	conn := NewIntegrationConnection(t)
	defer conn.Close(func(err error) {
		assert.NoError(t, err)
	})

	// "SELECT FROM pg_database WHERE datname = 'test'"
	rows, err := conn.Client().Query(
		context.Background(),
		"SELECT datname FROM pg_database WHERE datname = 'test'",
	)
	assert.NoError(t, err)
	for rows.Next() {
		var datname string
		err := rows.Scan(&datname)
		assert.NoError(t, err)
		fmt.Println(datname)
	}
	rows.Close()
	assert.Fail(t, "no")

	// // Create database "test" if non existent
	// _, err = conn.Client().Exec(
	// 	context.Background(),
	// 	"SELECT 'CREATE DATABASE test' WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'test')\\gexec",
	// )
	// assert.NoError(t, err)

	// Create a new database conn
	testDb, done, err := conn.UseDatabase("test")
	defer done()
	assert.NoError(t, err)

	// Create a new table
	_, err = testDb.Client().Exec(
		context.Background(),
		"CREATE TABLE IF NOT EXISTS test (id SERIAL PRIMARY KEY, name VARCHAR(255) NOT NULL)",
	)
	assert.NoError(t, err)
}
