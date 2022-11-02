package mongodb

import (
	"context"

	"github.com/Becklyn/clerk/v2"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type indexCreator struct {
	client     *mongo.Client
	collection *clerk.Collection
}

func newIndexCreator(connection *Connection, collection *clerk.Collection) *indexCreator {
	return &indexCreator{
		client:     connection.client,
		collection: collection,
	}
}

func (c *indexCreator) ExecuteCreate(
	ctx context.Context,
	create *clerk.Create[*clerk.Index],
) error {
	models := []mongo.IndexModel{}
	for _, index := range create.Data {
		fields := bson.D{}
		for _, field := range index.Fields {
			fields = append(fields, bson.E{
				Key: field.Key,
				Value: func() any {
					switch field.Type.String() {
					case "ascending":
						return 1
					case "descending":
						return -1
					case "text":
						return "text"
					}
					return nil
				}(),
			})
		}

		opts := options.
			Index().
			SetName(index.Name).
			SetUnique(index.IsUnique)

		models = append(models, mongo.IndexModel{
			Keys:    fields,
			Options: opts,
		})
	}

	_, err := c.client.
		Database(c.collection.Database.Name).
		Collection(c.collection.Name).
		Indexes().
		CreateMany(ctx, models)

	return err
}
