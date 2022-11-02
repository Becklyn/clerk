package clerk

type Database struct {
	Name string
}

func NewDatabase(name string) *Database {
	return &Database{
		Name: name,
	}
}
