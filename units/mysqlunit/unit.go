package mysqlunit

import (
	"github.com/dogukanayd/go-test-database/databases/testmysql"
	"testing"
	"time"

	"upper.io/db.v3/lib/sqlbuilder"
)

type MysqlInterface interface {
	Start() (sqlbuilder.Database, func())
}

type Unit struct {
	t            *testing.T
	databaseName string
}

// NewUnit generates a new unit instance
func NewUnit(t *testing.T) *Unit {
	return &Unit{t: t}
}

// NewUnit ...
func (u *Unit) Start() (sqlbuilder.Database, func()) {
	u.t.Helper()
	c := testmysql.NewContainer(u.t)

	connection, err := testmysql.Connections.ConnectOrReuse(u.databaseName, c.Host)

	if err != nil {
		u.t.Fatalf("opening database connection: %v", err)
	}

	u.HealthCheck(connection, u.t, c)

	// teardown is the function that should be invoked when the caller is done
	// with the database.
	teardown := func() {
		u.t.Helper()
		_ = connection.Close()
		c.StopContainer(u.t)
	}

	return connection, teardown
}

// Wait for the database to be ready. Wait 100ms longer between each attempt.
// Do not try more than 20 times.
func (u *Unit) HealthCheck(connection sqlbuilder.Database, t *testing.T, c *testmysql.Container) {
	var pingError error

	maxAttempts := 20

	for attempts := 1; attempts <= maxAttempts; attempts++ {
		pingError = connection.Ping()

		if pingError == nil {
			break
		}

		time.Sleep(time.Duration(attempts) * 100 * time.Millisecond)
	}

	if pingError != nil {
		c.DumpContainerLogs(t)
		c.StopContainer(t)
		t.Fatalf("waiting for database to be ready: %v", pingError)
	}
}
