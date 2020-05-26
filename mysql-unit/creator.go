package mysql_unit

import (
	"docker.io/go-docker"
	"docker.io/go-docker/api/types/container"
	"github.com/dogukanayd/go-test-database/containers"
)

type ContainerProperties struct {
	Config   *container.Config
	HostPort string
	CLI      *docker.Client
	TCP      string
	Name     string
}

func (cp *ContainerProperties) Creator(d containers.Docker) *containers.Container {
	d.SetConfig(cp.Config)
	d.SetContainerHostPort(cp.HostPort)
	d.SetTCP(cp.TCP)
	d.SetContainerCLI(cp.CLI)
	d.SetContainerName(cp.Name)

	return d.CreateContainer()
}
