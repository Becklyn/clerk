package clerk

type Collection struct {
	Database *Database
	Name     string
}

func NewCollection(database *Database, name string) *Collection {
	return &Collection{
		Database: database,
		Name:     name,
	}
}
