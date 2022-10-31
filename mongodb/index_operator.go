package mongodb

import (
	"context"

	"github.com/Becklyn/clerk"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongodbIndexOperator struct {
	client *mongo.Client
}

func NewMongoIndexOperator(connection *MongodbConnection) *MongodbIndexOperator {
	return &MongodbIndexOperator{
		client: connection.client,
	}
}

func (o *MongodbIndexOperator) List(
	ctx context.Context,
	collection *clerk.Collection,
) ([]*clerk.Index, error) {
	cursor, err := o.client.
		Database(collection.Database).
		Collection(collection.Name).
		Indexes().
		List(ctx)
	if err != nil {
		return nil, err
	}

	indices := []*clerk.Index{}
	models := []primitive.D{}
	if err := cursor.All(ctx, &models); err != nil {
		return nil, err
	}
	for _, indexModel := range models {
		index := &clerk.Index{
			Collection: collection,
		}
		for _, kv := range indexModel {
			switch kv.Key {
			case "name":
				index.Name = kv.Value.(string)
			case "key":
				indexFields := []*clerk.IndexField{}
				for _, field := range kv.Value.(primitive.D) {
					indexFields = append(indexFields, &clerk.IndexField{
						Key:  field.Key,
						Type: field.Value,
					})
				}
				index.Fields = indexFields
			case "unique":
				index.Unique = kv.Value.(bool)
			}
		}
		indices = append(indices, index)
	}
	return indices, nil
}

func (o *MongodbIndexOperator) Create(
	ctx context.Context,
	collection *clerk.Collection,
	indices ...*clerk.IndexCreate,
) ([]string, error) {
	models := []mongo.IndexModel{}
	for _, index := range indices {
		keys := bson.D{}
		for _, field := range index.Fields {
			keys = append(keys, bson.E{
				Key:   field.Key,
				Value: field.Type,
			})
		}

		options := options.
			Index().
			SetName(index.Name).
			SetUnique(index.Unique)

		model := mongo.IndexModel{
			Keys:    keys,
			Options: options,
		}

		models = append(models, model)
	}

	names, err := o.client.
		Database(collection.Database).
		Collection(collection.Name).
		Indexes().
		CreateMany(ctx, models)
	if err != nil {
		return nil, err
	}
	return names, nil
}

func (o *MongodbIndexOperator) Delete(
	ctx context.Context,
	collection *clerk.Collection,
	index *clerk.IndexDelete,
) error {
	_, err := o.client.
		Database(collection.Database).
		Collection(collection.Name).
		Indexes().
		DropOne(ctx, index.Name)

	return err
}
