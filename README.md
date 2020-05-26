# go-test-database

This package creates a mysql container for your applications and returns a connectionn to the created database so you can easily query the test database.

# Requirements
* :3305 port must not be used
* docker

# INSTALL
```bash
go get github.com/dogukanayd/go-test-database
``` 

# USAGE
```go
package mysql_unit

import (
	"log"
	"testing"
)

func TestNewUnit(t *testing.T) {
	connection, def := NewUnit()
    // example query
	q := `CREATE TABLE test_table(id int(11),name varchar(500)) ENGINE = InnoDB  DEFAULT CHARSET = utf8;`
	_, err := connection.Query(q)

	if err != nil {
		log.Fatal(err)
	}

	// TODO: connection should close

	defer def()

	t.Log("ping success")
}

```

When you call the `NewUnit` function;
```go
connection, tearDown := NewUnit()
```
it's return two parameters;
 * connection: *sql.DB // connection that allows you to query the database 
 * tearDown
 
here is the defination of the `NewUnit` function;
```go
func NewUnit() (*sql.DB, func()) {}
```
