# go-test-database
#### What this package does?
This package creates up and running mysql container(using docker sdk) for your applications and returns a connection
to the created database, so you can easily query the test database.

![alt text](https://encrypted-tbn0.gstatic.com/images?q=tbn%3AANd9GcS_5LAwzTqOQRs8pYq_fJHA8n7djUMdcgkeR2Qg69ajuggKXhgm&usqp=CAU)

# Requirements
* :3305 port must not be used
* docker

# INSTALL
```bash
go get github.com/dogukanayd/go-test-database
``` 

# USAGE
Let's say you have a package named greetings also in this package you have a function named `Hello`
and this function send some query to the database, well from at this point you need a database separated
from your local database. Here is all you need to do.

```go
connection, def := mysql_unit.NewUnit()
```

and below you can find full dummy example

```go
package greetings

import (
	mysql_unit "github.com/dogukanayd/go-test-database/mysql-unit"
	"log"
	"testing"
)

func TestHello(t *testing.T) {
	connection, def := mysql_unit.NewUnit()

	q := `CREATE TABLE test_table(id int(11),name varchar(500)) ENGINE = InnoDB  DEFAULT CHARSET = utf8;`
	_, err := connection.Query(q)
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
