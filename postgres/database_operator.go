package postgres

type DatabaseOperator struct {
	databaseCreator
	databaseQuerier
	databaseDeleter
	transactor
}

func NewDatabaseOperator(connection *Connection) *DatabaseOperator {
	return &DatabaseOperator{
		databaseCreator: *newDatabaseCreator(connection),
		databaseQuerier: *newDatabaseQuerier(connection),
		databaseDeleter: *newDatabaseDeleter(connection),
		transactor:      *newTransactor(),
	}
}
