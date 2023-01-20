package postgres

import sq "github.com/Masterminds/squirrel"

func statementBuilder() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}
