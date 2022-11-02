package mongodb

type databaseOperator struct {
	databaseQuerier
	transactor
}

func NewDatabaseOperator(connection *Connection) *databaseOperator {
	return &databaseOperator{
		databaseQuerier: *newDatabaseQuerier(connection),
		transactor:      *newTransactor(connection),
	}
}
