package mysql

import (
	"docker.io/go-docker"
	"docker.io/go-docker/api/types/container"
)

type ContainerProperties struct {
	Config   *container.Config
	HostPort string
	CLI      *docker.Client
	TCP      string
	Name     string
}

func (cp *ContainerProperties) Creator(d Docker) *DockerContainer {
	d.SetConfig(cp.Config)
	d.SetContainerHostPort(cp.HostPort)
	d.SetTCP(cp.TCP)
	d.SetContainerCLI(cp.CLI)
	d.SetContainerName(cp.Name)

	return d.CreateContainer()
}
