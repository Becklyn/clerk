package mongodb

import (
	"context"

	"github.com/Becklyn/clerk/v2"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type indexQuerier struct {
	client     *mongo.Client
	collection *clerk.Collection
}

func newIndexQuerier(connection *Connection, collection *clerk.Collection) *indexQuerier {
	return &indexQuerier{
		client:     connection.client,
		collection: collection,
	}
}

func (q *indexQuerier) ExecuteQuery(
	ctx context.Context,
	query *clerk.Query[*clerk.Index],
) (<-chan *clerk.Index, error) {
	cursor, err := q.client.
		Database(q.collection.Database.Name).
		Collection(q.collection.Name).
		Indexes().
		List(ctx)
	if err != nil {
		return nil, err
	}

	channel := make(chan *clerk.Index)

	go func() {
		defer cursor.Close(ctx)
		defer close(channel)

		for cursor.Next(ctx) {
			model := primitive.D{}
			if err := cursor.Decode(&model); err != nil {
				return
			}

			indx := &clerk.Index{}
			for _, kv := range model {
				switch kv.Key {
				case "name":
					indx.Name = kv.Value.(string)
				case "key":
					fields := []*clerk.Field{}
					for _, field := range kv.Value.(primitive.D) {
						fields = append(fields, &clerk.Field{
							Key: field.Key,
							Type: func() clerk.FieldType {
								switch field.Value {
								case 1:
									return clerk.FieldTypeAscending
								case -1:
									return clerk.FieldTypeDescending
								case "text":
									return clerk.FieldTypeString
								}
								return clerk.FieldTypeAscending
							}(),
						})
					}
					indx.Fields = fields
				case "unique":
					indx.IsUnique = kv.Value.(bool)
				}
			}

			if indx.Name == "_id_" {
				indx.IsUnique = true
			}

			channel <- indx
		}
	}()

	return channel, nil
}
