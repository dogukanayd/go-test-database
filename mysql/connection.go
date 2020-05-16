package mysql

import (
	"time"
	"upper.io/db.v3/lib/sqlbuilder"
	"upper.io/db.v3/mysql"
)

// Connections
var Connections = &connection{list: make(map[string]sqlbuilder.Database)}

type connection struct {
	list map[string]sqlbuilder.Database
}

const TestDatabaseName = "test_database"

// ConnectOrReuse connection
func (c *connection) ConnectOrReuse(dbName, host string) (sqlbuilder.Database, error) {
	connection, ok := c.list[dbName]
	var session sqlbuilder.Database
	var sessionError error

	// If connection does not exist or dead
	if !ok || connection.Ping() != nil {
		url := "root" + ":" + "root" + "@" + "tcp(" + host + ")"
		url = url + "/" + dbName + "?charset=utf8mb4&collation=utf8mb4_unicode_ci"
		dsn, err := mysql.ParseURL(url)

		if err != nil {
			return nil, err
		}

		maxAttempts := 20

		for attempts := 1; attempts <= maxAttempts; attempts++ {
			session, sessionError = mysql.Open(dsn)

			if sessionError == nil {
				break
			}

			time.Sleep(time.Duration(attempts) * 1000 * time.Millisecond)
		}

		c.list[dbName] = session

		return session, nil
	}

	return connection, nil
}
