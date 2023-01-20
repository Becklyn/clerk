package postgres

type DatabaseOperator struct {
	databaseQuerier
}

func NewDatabaseOperator(connection *Connection) *DatabaseOperator {
	return &DatabaseOperator{
		databaseQuerier: *newDatabaseQuerier(connection),
	}
}
