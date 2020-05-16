package mysql

import (
	"testing"
	"time"

	"upper.io/db.v3/lib/sqlbuilder"
)

type Unit struct {
	T            *testing.T
	DatabaseName string
	MaxAttempts  int
}

// NewUnit generates a new unit instance
//
// Example:
// 		mysqlunit.NewUnit(t, "test_database").Start()
func NewUnit(t *testing.T, dbName string, ma int) *Unit {
	return &Unit{
		T:            t,
		DatabaseName: dbName,
		MaxAttempts:  ma,
	}
}

// NewUnit ...
func (u *Unit) Start() (sqlbuilder.Database, func()) {
	u.T.Helper()
	c := NewContainer()

	connection, err := Connections.ConnectOrReuse(u.DatabaseName, c.Host)

	if err != nil {
		u.T.Fatalf("opening database connection: %v", err)
	}

	u.HealthCheck(connection, u.T, c)

	// teardown is the function that should be invoked when the caller is done
	// with the database.
	teardown := func() {
		u.T.Helper()
		_ = connection.Close()
		c.StopContainer(u.T)
	}

	return connection, teardown
}

// Wait for the database to be ready. Wait 100ms longer between each attempt.
// Do not try more than 20 times.
func (u *Unit) HealthCheck(connection sqlbuilder.Database, t *testing.T, c *Container) {
	var pingError error

	maxAttempts := 20

	for attempts := 1; attempts <= maxAttempts; attempts++ {
		pingError = connection.Ping()

		if pingError == nil {
			break
		}

		time.Sleep(time.Duration(attempts) * 1000 * time.Millisecond)
	}

	if pingError != nil {
		c.DumpContainerLogs(t)
		c.StopContainer(t)
		t.Fatalf("waiting for database to be ready: %v", pingError)
	}
}
