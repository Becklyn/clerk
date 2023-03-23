package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/Becklyn/clerk/v4"
	"github.com/jackc/pgx/v5"
	"github.com/samber/lo"
)

type indexQuerier struct {
	conn              *Connection
	collection        *clerk.Collection
	collectionCreator *collectionCreator
	transactor        *transactor
}

func newIndexQuerier(conn *Connection, collection *clerk.Collection) *indexQuerier {
	return &indexQuerier{
		conn:              conn,
		collection:        collection,
		collectionCreator: newCollectionCreator(conn, collection.Database),
		transactor:        newTransactor(conn),
	}
}

func (q *indexQuerier) Count(
	ctx context.Context,
	query *clerk.Query[*clerk.Index],
) (int64, error) {
	var name string
	for _, filter := range query.Filters {
		switch filter.(type) {
		case *clerk.Equals:
			if strings.ToLower(filter.Key()) == "name" {
				name = fmt.Sprintf("%s_%s", q.collection.Name, filter.Value().(string))
			}
		}
	}

	selectFn := func(ctx context.Context, dbConn *pgx.Conn) (pgx.Rows, error) {
		if name == "" {
			return dbConn.Query(ctx, "SELECT Count(*) FROM pg_indexes WHERE tablename = $1", q.collection.Name)
		}

		return dbConn.Query(ctx, "SELECT Count(*) FROM pg_indexes WHERE tablename = $1 AND indexname = $2", q.collection.Name, name)
	}

	queryCtx, cancel := q.conn.config.GetContext(ctx)
	defer cancel()

	var total int64

	if err := q.transactor.executeInTransactionIfAvailable(queryCtx, q.collection.Database, func(ctx context.Context) error {
		dbConn, release, err := q.conn.createOrUseDatabase(ctx, q.collection.Database.Name)
		defer release()
		if err != nil {
			return err
		}

		rows, err := selectFn(ctx, dbConn)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			if err := rows.Scan(&total); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return 0, err
	}

	return total, nil
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

	selectFn := func(ctx context.Context, dbConn *pgx.Conn) (pgx.Rows, error) {
		if name == "" {
			return dbConn.Query(ctx, "SELECT indexname, indexdef FROM pg_indexes WHERE tablename = $1", q.collection.Name)
		}

		return dbConn.Query(ctx, "SELECT indexname, indexdef FROM pg_indexes WHERE tablename = $1 AND indexname = $2", q.collection.Name, name)
	}

	queryCtx, cancel := q.conn.config.GetContext(ctx)
	defer cancel()

	var indices []*clerk.Index

	if err := q.transactor.executeInTransactionIfAvailable(queryCtx, q.collection.Database, func(ctx context.Context) error {
		dbConn, release, err := q.conn.createOrUseDatabase(ctx, q.collection.Database.Name)
		defer release()
		if err != nil {
			return err
		}

		rows, err := selectFn(ctx, dbConn)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			index := &clerk.Index{}

			var (
				indexName string
				indexDef  string
			)

			if err := rows.Scan(&indexName, &indexDef); err != nil {
				return err
			}

			index.Name = strings.TrimPrefix(indexName, q.collection.Name+"_")
			index.IsUnique = strings.HasPrefix(indexDef, "CREATE UNIQUE INDEX ")
			index.Fields = lo.Map(getFieldDefs(indexDef), func(fieldDef string, _ int) *clerk.Field {
				return getFieldFromDef(fieldDef)
			})

			indices = append(indices, index)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	channel := make(chan *clerk.Index)

	go func() {
		defer close(channel)

		for _, index := range indices {
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
