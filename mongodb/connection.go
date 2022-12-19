package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Connection struct {
	ctx    context.Context
	config Config
	client *mongo.Client
}

func NewConnection(
	ctx context.Context,
	config Config,
) (*Connection, error) {
	opts := options.Client().
		SetConnectTimeout(config.Timeout).
		SetServerSelectionTimeout(config.Timeout).
		SetSocketTimeout(config.Timeout).
		SetTimeout(config.Timeout).
		ApplyURI(config.Uri)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	pingCtx, pingCancel := config.GetContext(ctx)
	defer pingCancel()
	if err = client.Ping(pingCtx, readpref.Primary()); err != nil {
		return nil, err
	}

	return &Connection{
		ctx:    ctx,
		config: config,
		client: client,
	}, nil
}

func (c *Connection) Close(handler func(err error)) {
	err := c.client.Disconnect(c.ctx)
	if handler != nil {
		handler(err)
	}
}
