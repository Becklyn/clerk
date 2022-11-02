package mongodb

import (
	"context"
	"strings"

	"github.com/Becklyn/clerk/v2"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type collectionUpdater struct {
	connection *Connection
	client     *mongo.Client
	database   *clerk.Database
}

func newCollectionUpdater(connection *Connection, database *clerk.Database) *collectionUpdater {
	return &collectionUpdater{
		connection: connection,
		client:     connection.client,
		database:   database,
	}
}

func (u *collectionUpdater) ExecuteUpdate(
	ctx context.Context,
	update *clerk.Update[*clerk.Collection],
) error {
	names := []string{}
	for _, filter := range update.Filters {
		switch filter.(type) {
		case *clerk.Equals:
			if strings.ToLower(filter.Key()) == "name" {
				names = append(names, filter.Value().(string))
			}
		}
	}

	err := newTransactor(u.connection).ExecuteTransaction(ctx, func(ctx context.Context) error {
		for _, name := range names {
			cursor, err := u.client.
				Database(u.database.Name).
				Collection(name).
				Find(ctx, bson.D{}, options.Find())
			if err != nil {
				return err
			}

			for cursor.Next(ctx) {
				var result any
				if err := cursor.Decode(&result); err != nil {
					return err
				}

				_, err = u.client.
					Database(u.database.Name).
					Collection(update.Data.Name).
					InsertOne(ctx, result)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	for _, name := range names {
		err := u.client.
			Database(u.database.Name).
			Collection(name).
			Drop(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
