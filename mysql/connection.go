package mysql

import (
	"fmt"
	"sync"
	"time"
	"upper.io/db.v3/lib/sqlbuilder"
	"upper.io/db.v3/mysql"
)

// Connections
var Connections = &connection{list: make(map[string]sqlbuilder.Database)}

type connection struct {
	list map[string]sqlbuilder.Database
	sync.Mutex
}

const TestDatabaseName = "test_database"

// ConnectOrReuse connection
func (c *connection) ConnectOrReuse(dbName, host string) (sqlbuilder.Database, error) {
	c.Lock()
	defer c.Unlock()

	connection, ok := c.list[dbName]

	// If connection does not exist or dead
	if !ok || connection.Ping() != nil {
		url := "root" + ":" + "root" + "@" + "tcp(" + host + ")"
		url = url + "/" + dbName + "?charset=utf8mb4&collation=utf8mb4_unicode_ci"
		dsn, err := mysql.ParseURL(url)

		if err != nil {
			return nil, err
		}

		session, err := mysql.Open(dsn)

		c.list[dbName] = session
		go c.close(dbName)

		return session, nil
	}

	return connection, nil
}

func (c *connection) close(db string) {
	time.Sleep(30 * time.Second)

	c.Lock()
	defer c.Unlock()

	if err := c.list[db].Close(); err != nil {
		fmt.Println(err)
	}

	delete(c.list, db)
}
