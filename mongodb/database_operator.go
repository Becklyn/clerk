package mongodb

type DatabaseOperator struct {
	databaseQuerier
	transactor
}

func NewDatabaseOperator(connection *Connection) *DatabaseOperator {
	return &DatabaseOperator{
		databaseQuerier: *newDatabaseQuerier(connection),
		transactor:      *newTransactor(connection),
	}
}
