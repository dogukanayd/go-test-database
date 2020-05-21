package mysql

import (
	"docker.io/go-docker"
	"docker.io/go-docker/api/types/container"
	"fmt"
	"upper.io/db.v3/lib/sqlbuilder"
)

func NewUnitV2() (sqlbuilder.Database, func()) {
	var cp ContainerProperties
	var err error

	cp.CLI, err = docker.NewEnvClient()

	if err != nil {
		panic(err)
	}

	cp.Config = &container.Config{
		Image: "mysql:5.6",
		Env: []string{
			"MYSQL_ROOT_PASSWORD=root",
			"MYSQL_DATABASE=test_database_v2",
		},
	}
	cp.HostPort = "3305"
	cp.TCP = "3306"
	cpr := cp.Creator(&DockerContainer{})

	fmt.Println("container id ******************")
	fmt.Println(cpr.ID)

	fmt.Println("connection start ******************")
	connection, err := Connections.ConnectOrReuse("test_database_v2", "127.0.0.1:3305")

	fmt.Println("connection end ******************")

	if err != nil {
		panic(err)
	}

	tearDown := func() {
		_ = connection
		_ = connection.Close()
	}

	return connection, tearDown
}
