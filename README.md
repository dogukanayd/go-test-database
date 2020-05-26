# go-test-database

This package creates a up and running mysql container for your applications and returns a connectionn to the created database so you can easily query the test database.

# Requirements
* :3305 port must not be used
* docker

# INSTALL
```bash
go get github.com/dogukanayd/go-test-database
``` 

# USAGE
```go
package greetings

import (
	mysql_unit "github.com/dogukanayd/go-test-database/mysql-unit"
	"log"
	"testing"
)

func TestHello(t *testing.T) {
	db, def := mysql_unit.NewUnit()

	q := `CREATE TABLE test_table(id int(11),name varchar(500)) ENGINE = InnoDB  DEFAULT CHARSET = utf8;`
	_, err := db.Query(q)
	if err != nil {
		log.Fatal(err)
	}

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
