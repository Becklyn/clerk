package postgres

type DatabaseOperator struct {
	databaseCreator
	databaseQuerier
	databaseDeleter
}

func NewDatabaseOperator(connection *Connection) *DatabaseOperator {
	return &DatabaseOperator{
		databaseCreator: *newDatabaseCreator(connection),
		databaseQuerier: *newDatabaseQuerier(connection),
		databaseDeleter: *newDatabaseDeleter(connection),
	}
}
