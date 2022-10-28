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
	client     *mongo.Client
	collection *clerk.Collection
}

func NewMongoIndexOperator(client *mongo.Client, collection *clerk.Collection) *MongodbIndexOperator {
	return &MongodbIndexOperator{
		client:     client,
		collection: collection,
	}
}

func (o *MongodbIndexOperator) List(ctx context.Context) ([]*clerk.Index, error) {
	cursor, err := o.client.
		Database(o.collection.Database).
		Collection(o.collection.Name).
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
		index := &clerk.Index{}
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

func (o *MongodbIndexOperator) Create(ctx context.Context, index *clerk.Index) error {
	options := options.Index()
	if index.Name != "" {
		options.SetName(index.Name)
	}
	if index.Unique {
		options.SetUnique(true)
	}

	modelKeys := bson.D{}
	for _, field := range index.Fields {
		modelKeys = append(modelKeys, bson.E{
			Key:   field.Key,
			Value: field.Type,
		})
	}

	model := mongo.IndexModel{
		Keys:    modelKeys,
		Options: options,
	}

	_, err := o.client.
		Database(o.collection.Database).
		Collection(o.collection.Name).
		Indexes().
		CreateOne(ctx, model)

	return err
}
