package postgres

import (
	"fmt"

	"github.com/jackc/pgx/v5"
)

type DatabaseConnection struct {
	conn   *Connection
	client *pgx.Conn
}

func NewDatabaseConnection(
	conn *Connection,
	database string,
) (*DatabaseConnection, error) {
	client, err := pgx.Connect(conn.ctx, fmt.Sprintf("%s/%s", conn.config.Host, database))
	if err != nil {
		return nil, err
	}

	return &DatabaseConnection{
		conn:   conn,
		client: client,
	}, nil
}

// TODO: DONT EXPORT THIS
func (c *DatabaseConnection) Client() *pgx.Conn {
	return c.client
}

func (c *DatabaseConnection) Close() error {
	return c.client.Close(c.conn.ctx)
}
