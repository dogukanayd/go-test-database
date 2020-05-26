# go-test-database

This package creates a test database for your go applications using the `127.0.0.1:3305`.

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
	db, def := NewUnit()
    // example query
	q := `CREATE TABLE test_table(id int(11),name varchar(500)) ENGINE = InnoDB  DEFAULT CHARSET = utf8;`
	_, err := db.Query(q)

	if err != nil {
		log.Fatal(err)
	}

	// TODO: connection should close

	defer def()

	t.Log("ping success")
}

```