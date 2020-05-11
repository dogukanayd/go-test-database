package mysqlunit

import (
	"github.com/dogukanayd/go-test-database/databases/testmysql"
	"testing"
	"time"

	"upper.io/db.v3/lib/sqlbuilder"
)

// NewUnit ...
func NewUnit(t *testing.T) (sqlbuilder.Database, func()) {
	t.Helper()
	c := testmysql.StartContainer(t)

	connection, err := testmysql.Connections.ConnectOrReuse(testmysql.TestDatabaseName, c.Host)

	if err != nil {
		t.Fatalf("opening database connection: %v", err)
	}

	healthCheck(connection, t, c)

	// teardown is the function that should be invoked when the caller is done
	// with the database.
	teardown := func() {
		t.Helper()
		_ = connection.Close()
		testmysql.StopContainer(t, c)
	}

	return connection, teardown
}

// Wait for the database to be ready. Wait 100ms longer between each attempt.
// Do not try more than 20 times.
func healthCheck(connection sqlbuilder.Database, t *testing.T, c *testmysql.Container, ) {
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
		testmysql.StopContainer(t, c)
		t.Fatalf("waiting for database to be ready: %v", pingError)
	}
}
