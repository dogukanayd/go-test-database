package mysql

import (
	"docker.io/go-docker"
	"docker.io/go-docker/api/types/container"
	"testing"
	"time"
	"upper.io/db.v3/lib/sqlbuilder"
)

type Unit struct {
	T            *testing.T
	DatabaseName string
}

// NewUnit generates a new unit instance
//
// Example:
// 		mysqlunit.NewUnit(t, "test_database").Start()
func NewUnit(t *testing.T, dbName string) *Unit {
	return &Unit{
		T:            t,
		DatabaseName: dbName,
	}
}

// NewUnit ...
func (u *Unit) Start() (sqlbuilder.Database, func()) {
	cli, err := docker.NewEnvClient()

	if err != nil {
		panic(err)
	}

	u.T.Helper()
	c := NewContainer(
		"3305",
		"3306",
		"",
		&container.Config{
			Image: "mysql:5.6",
			Env: []string{
				"MYSQL_ROOT_PASSWORD=root",
				"MYSQL_DATABASE=test_database_v2",
			},
		},
		cli,
	)
	c.StartContainer()

	retry := 20

	for i := 1; i <= retry; i++ {
		if status, err := c.GetStatus(); status[0].State.Status == "running" &&
			status[0].NetworkSettings.Networks.Bridge.IPAddress != "" &&
			err == nil {
			break
		}

		time.Sleep(2 * time.Second)
	}

	connection, err := Connections.ConnectOrReuse(u.DatabaseName, "127.0.0.1:3305")

	if err != nil {
		u.T.Fatalf("opening database connection: %v", err)
	}

	// teardown is the function that should be invoked when the caller is done
	// with the database.
	teardown := func() {
		u.T.Helper()
		_ = connection.Close()
		c.StopContainer()
	}

	return connection, teardown
}
