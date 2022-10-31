package mongodb

import (
	"context"

	"github.com/Becklyn/clerk"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongodbCollectionOperator struct {
	connection *MongodbConnection
	client     *mongo.Client
}

func NewMongoCollectionOperator(connection *MongodbConnection) *MongodbCollectionOperator {
	return &MongodbCollectionOperator{
		connection: connection,
		client:     connection.client,
	}
}

func (o *MongodbCollectionOperator) List(
	ctx context.Context,
	database *clerk.Database,
) ([]*clerk.Collection, error) {
	cursor, err := o.client.
		Database(database.Name).
		ListCollectionNames(ctx, bson.D{}, options.ListCollections())
	if err != nil {
		return nil, err
	}

	var collections []*clerk.Collection
	for _, name := range cursor {
		collections = append(collections, clerk.NewCollectionWithDatabase(database.Name, name))
	}
	return collections, nil
}

func (o *MongodbCollectionOperator) Rename(
	ctx context.Context,
	collection *clerk.Collection,
	renameTo string,
	drop bool,
) (*clerk.Collection, error) {
	handler := func(ctx context.Context) error {
		cursor, err := o.client.
			Database(collection.Database).
			Collection(collection.Name).
			Find(ctx, bson.D{}, options.Find())
		if err != nil {
			return err
		}
		var results []interface{}
		if err := cursor.All(ctx, &results); err != nil {
			return err
		}

		_, err = o.client.Database(collection.Database).
			Collection(renameTo).
			InsertMany(ctx, results)
		if err != nil {
			return err
		}
		return nil
	}

	if err := o.connection.WithTransaction(ctx, handler); err != nil {
		return nil, err
	}
	if drop {
		if err := o.Drop(ctx, collection); err != nil {
			return nil, err
		}
	}
	return clerk.NewCollectionWithDatabase(collection.Database, renameTo), nil
}

func (o *MongodbCollectionOperator) Drop(
	ctx context.Context,
	collection *clerk.Collection,
) error {
	return o.client.
		Database(collection.Database).
		Collection(collection.Name).
		Drop(ctx)
}
