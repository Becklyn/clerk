package mongodb

type DatabaseOperator struct {
	databaseQuerier
	databaseDeleter
	transactor
}

func NewDatabaseOperator(connection *Connection) *DatabaseOperator {
	return &DatabaseOperator{
		databaseQuerier: *newDatabaseQuerier(connection),
		databaseDeleter: *newDatabaseDeleter(connection),
		transactor:      *newTransactor(connection),
	}
}
