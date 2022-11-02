package mongodb

import (
	"context"
	"strings"

	"github.com/Becklyn/clerk/v2"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var operationMap = map[clerk.Operation]string{
	clerk.OperationCreate: "insert",
	clerk.OperationDelete: "delete",
	clerk.OperationUpdate: "replace|update",
}

type watcher[T any] struct {
	client     *mongo.Client
	collection *clerk.Collection
}

func newWatcher[T any](connection *Connection, collection *clerk.Collection) *watcher[T] {
	return &watcher[T]{
		client:     connection.client,
		collection: collection,
	}
}

func (w *watcher[T]) ExecuteWatch(
	ctx context.Context,
	watch *clerk.Watch[T],
) (<-chan *clerk.Event[T], error) {
	opts := options.ChangeStream().
		SetFullDocument(options.UpdateLookup)

	stream, err := w.client.
		Database(w.collection.Database.Name).
		Collection(w.collection.Name).
		Watch(ctx, mongo.Pipeline{}, opts)
	if err != nil {
		return nil, err
	}

	channel := make(chan *clerk.Event[T])

	go func() {
		defer close(channel)
		defer stream.Close(ctx)

		for stream.Next(ctx) {
			var event bson.M
			if err := stream.Decode(&event); err != nil {
				continue
			}

			operationType, ok := event["operationType"].(string)
			if !ok {
				continue
			}
			operation, ok := func() (clerk.Operation, bool) {
				for op, pattern := range operationMap {
					if strings.Contains(pattern, operationType) {
						return op, true
					}
				}
				return -1, false
			}()
			if !ok {
				continue
			}

			document, ok := event["documentKey"].(bson.M)
			if !ok {
				continue
			}
			if operationType != "delete" {
				document, ok = event["fullDocument"].(bson.M)
				if !ok {
					continue
				}
			}

			documentData, err := bson.Marshal(document)
			if err != nil {
				continue
			}

			var data T
			if err := bson.Unmarshal(documentData, &data); err != nil {
				continue
			}

			channel <- &clerk.Event[T]{
				Operation: operation,
				Data:      data,
			}
		}
	}()

	return channel, nil
}
