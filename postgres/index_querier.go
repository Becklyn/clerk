package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/Becklyn/clerk/v3"
	"github.com/jackc/pgx/v5"
	"github.com/samber/lo"
)

type indexQuerier struct {
	conn       *Connection
	collection *clerk.Collection
}

func newIndexQuerier(conn *Connection, collection *clerk.Collection) *indexQuerier {
	return &indexQuerier{
		conn:       conn,
		collection: collection,
	}
}

func (q *indexQuerier) ExecuteQuery(
	ctx context.Context,
	query *clerk.Query[*clerk.Index],
) (<-chan *clerk.Index, error) {
	var name string
	for _, filter := range query.Filters {
		switch filter.(type) {
		case *clerk.Equals:
			if strings.ToLower(filter.Key()) == "name" {
				name = fmt.Sprintf("%s_%s", q.collection.Name, filter.Value().(string))
			}
		}
	}

	queryCtx, cancel := q.conn.config.GetContext(ctx)

	dbConn, release, err := q.conn.useDatabase(queryCtx, q.collection.Database.Name)
	if err != nil {
		release()
		cancel()
		return nil, err
	}

	rows, err := func() (pgx.Rows, error) {
		if name == "" {
			return dbConn.Query(ctx, "SELECT indexname, indexdef FROM pg_indexes WHERE tablename = $1", q.collection.Name)
		}

		return dbConn.Query(ctx, "SELECT indexname, indexdef FROM pg_indexes WHERE tablename = $1 AND indexname = $2", q.collection.Name, name)
	}()
	if err != nil {
		release()
		cancel()
		return nil, err
	}

	channel := make(chan *clerk.Index)

	go func() {
		defer rows.Close()
		defer release()
		defer cancel()
		defer close(channel)

		for rows.Next() {
			index := &clerk.Index{}

			var (
				indexName string
				indexDef  string
			)

			if err := rows.Scan(&indexName, &indexDef); err != nil {
				return
			}

			index.Name = strings.TrimPrefix(indexName, q.collection.Name+"_")
			index.IsUnique = strings.HasPrefix(indexDef, "CREATE UNIQUE INDEX ")
			index.Fields = lo.Map(getFieldDefs(indexDef), func(fieldDef string, _ int) *clerk.Field {
				return getFieldFromDef(fieldDef)
			})

			channel <- index
		}
	}()

	return channel, nil
}

func getFieldDefs(indexDef string) []string {
	columnsStart := strings.Index(indexDef, "(")
	allColumnsString := indexDef[columnsStart+1 : len(indexDef)-1]
	return strings.Split(allColumnsString, ", ")
}

func getFieldFromDef(fieldDef string) *clerk.Field {
	field := &clerk.Field{}

	if strings.HasSuffix(fieldDef, ") DESC") {
		field.Type = clerk.FieldTypeDescending
	} else {
		field.Type = clerk.FieldTypeAscending
	}

	lastFieldIndex := strings.Index(fieldDef, " ->")
	field.Key = fieldDef[2:lastFieldIndex]

	return field
}
