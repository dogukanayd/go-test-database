package testmysql

import (
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

	// If connection does not exist or dead
	if !ok || connection.Ping() != nil {
		url := "root" + ":" + "root" + "@" + "tcp(" + host + ")"
		url = url + "/" + dbName + "?charset=utf8mb4&collation=utf8mb4_unicode_ci"
		dsn, err := mysql.ParseURL(url)

		if err != nil {
			return nil, err
		}

		session, err := mysql.Open(dsn)

		if err != nil {
			return nil, err
		}

		c.list[dbName] = session

		return session, nil
	}

	return connection, nil
}
