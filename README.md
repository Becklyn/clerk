# clerk

📒 A minimalistic library for abstracting database operations

## Installation

Adding *clerk* to your Go module is as easy as calling this command in your project

```shell
go get github.com/Becklyn/clerk/v2
```

## Supported databases

*clerk* has builtin support for the following database engines:

- [MongoDB](https://www.mongodb.com/) - MongoDB is a document-oriented database

Support for any other engines can be added by implementing their supported operations.  
Have a look at the *MongoDB* implementation in the `/mongodb` directory as a starting point.

## Usage

Being a minimalistic library, _clerk_ only provides the basics. The rest is up to your specific need.  
Each operation in *clerk* consists of a generic operation and uses an *operator* that is specific to the database engine.  
The examples given below all use the *MongoDB* operators.

### Creating a connection

```go
connection, err := mongodb.NewConnection(context.Background(), "mongodb://localhost:27017")
if err != nil {
	panic(err)
}

defer connection.Close(func(err error) {
	if err != nil {
		panic(err)
	}
})
```

### Using a transaction

```go
databaseOperator := mongodb.NewDatabaseOperator(connection)

clerk.NewTransaction(databaseOperator).Run(context.Background(), func(ctx context.Context) error {
    // Add operations that should be executed in a transaction here ...
    return nil
})
```

### Defining a database & collection

```go
database := clerk.NewDatabase("foo")
collection := clerk.NewCollection(database, "bar")
```

### Defining a database operator

```go
tOperator := mongodb.NewOperator[T](connection, collection)
```

The generic parameter T defines the data type which the operator can interact with.
An operator has to be defined for each data type in use with *clerk*.

### Persisting new data in a collection

```go
type Message struct {
    Id   string `bson:"_id"`
    Body string `bson:"body"`
}

messageOperator := mongodb.NewOperator[*Message](connection, collection)

createCtx, createCancel := context.WithTimeout(context.Background(), 5*time.Second)
defer createCancel()

err := clerk.NewCreate[*Message](messageOperator).
    With(&Message{Id: 1, Body: "Hello World"}).
    With(&Message{Id: 2, Body: "Hello Buddy"}).
    Commit(createCtx)
if err != nil {
    panic(err)
}
```

### Querying the collection

```go
type Message struct {
    Id   string `bson:"_id"`
    Body string `bson:"body"`
}

messageOperator := mongodb.NewOperator[*Message](connection, collection)

queryCtx, queryCancel := context.WithTimeout(context.Background(), 5*time.Second)
defer queryCancel()

message, err := clerk.NewQuery[*Message](messageOperator).
    Where(clerk.NewEquals("_id", 1)).
    Single(queryCtx)
if err != nil {
    panic(err)
}

fmt.Printf("Message: %+v", message)
```

```go
messageOperator := mongodb.NewOperator[*Message](connection, collection)

queryCtx, queryCancel := context.WithTimeout(context.Background(), 5*time.Second)
defer queryCancel()

messages, err := clerk.NewQuery[*Message](messageOperator).
    Where(clerk.NewRegex("body", "^Hello.*$")).
    Sort("_id", clerk.NewAscendingOrder()).
    All(queryCtx)
if err != nil {
    panic(err)
}

for _, message := range messages {
    fmt.Printf("Message: %+v", message)
}
```

---

Copyright © 2022 **Becklyn GmbH**
