package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Becklyn/clerk/v4"
	"github.com/samber/lo"
)

var ErrTextIndexNotSupported = errors.New("cannot use text indices with postgres")

type indexCreator struct {
	conn              *Connection
	collection        *clerk.Collection
	collectionCreator *collectionCreator
}

func newIndexCreator(conn *Connection, collection *clerk.Collection) *indexCreator {
	return &indexCreator{
		conn:              conn,
		collection:        collection,
		collectionCreator: newCollectionCreator(conn, collection.Database),
	}
}

func (c *indexCreator) ExecuteCreate(
	ctx context.Context,
	create *clerk.Create[*clerk.Index],
) error {
	createCtx, cancel := c.conn.config.GetContext(ctx)
	defer cancel()

	dbConn, release, err := c.conn.useDatabase(createCtx, c.collection.Database.Name)
	defer release()
	if err != nil {
		return err
	}

	for _, index := range create.Data {
		for _, field := range index.Fields {
			if field.Type.String() == "text" {
				return ErrTextIndexNotSupported
			}
		}

		var unique string

		if index.IsUnique {
			unique = "UNIQUE "
		}

		indexName := fmt.Sprintf("%s_%s", c.collection.Name, index.Name)

		columns := lo.Reduce(index.Fields, func(columnString string, field *clerk.Field, i int) string {
			order := func() string {
				if field.Type.String() == "descending" {
					return "DESC"
				}

				return "ASC"
			}()

			fieldName := jsonKeyToSelector("data", field.Key, nil)

			if i == 0 {
				return fmt.Sprintf("(%s) %s", fieldName, order)
			}
			return fmt.Sprintf("%s, (%s) %s", columnString, fieldName, order)
		}, "")

		stmt := fmt.Sprintf("CREATE %sINDEX IF NOT EXISTS %s ON %s (%s)", unique, indexName, c.collection.Name, columns)

		if _, err = dbConn.Exec(createCtx, stmt); err != nil {
			if err := c.collectionCreator.ExecuteCreate(ctx, &clerk.Create[*clerk.Collection]{
				Data: []*clerk.Collection{
					c.collection,
				},
			}); err != nil {
				return err
			}

			if _, err = dbConn.Exec(createCtx, stmt); err != nil {
				return err
			}
		}
	}

	return nil
}
