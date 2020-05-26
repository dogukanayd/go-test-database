package mysql_unit

import (
	"database/sql"
	"docker.io/go-docker"
	"docker.io/go-docker/api/types/container"
	"fmt"
	"github.com/dogukanayd/go-test-database/containers"
	"time"
)

func NewUnit() (*sql.DB, func()) {
	var cp ContainerProperties
	var connection *sql.DB
	var err error

	cp.CLI, err = docker.NewEnvClient()

	if err != nil {
		panic(err)
	}

	cp.Config = &container.Config{
		Image: "mysql:5.6",
		Env: []string{
			"MYSQL_ROOT_PASSWORD=root",
			"MYSQL_DATABASE=test_database",
		},
	}
	cp.HostPort = "3305"
	cp.TCP = "3306"

	cp.Creator(&containers.Container{})

	maxAttempts := 40

	connection = BasicConnection()

	for attempts := 1; attempts <= maxAttempts; attempts++ {
		fmt.Printf("Mysql connection attempt number: %v\n", attempts)

		if err = connection.Ping(); err != nil {
			// when ping failed error printed
			fmt.Println("trying to reconnect to the mysql...")

			time.Sleep(time.Duration(attempts) * 5 * time.Second)
		} else {
			break
		}
	}

	if err != nil {
		panic(err)
	}

	tearDown := func() {
		_ = connection
		_ = connection.Close()
	}

	return connection, tearDown
}
