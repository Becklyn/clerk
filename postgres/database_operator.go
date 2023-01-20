package postgres

type DatabaseOperator struct {
	databaseCreator
	databaseQuerier
}

func NewDatabaseOperator(connection *Connection) *DatabaseOperator {
	return &DatabaseOperator{
		databaseCreator: *newDatabaseCreator(connection),
		databaseQuerier: *newDatabaseQuerier(connection),
	}
}
