package mysql

import (
	"context"
	"docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/container"
	"github.com/docker/go-connections/nat"
)

type Docker interface {
	SetContainerID(id string)
	SetContainerName(name string)
	SetContainerHost(host string)
	SetContainerHostIP(hostIp string)
	SetContainerHostPort(hostPort string)
	SetTCP(tcp string)
	SetConfig(config *container.Config)
	SetContainerCLI(client *docker.Client)
	CreateContainer() *DockerContainer
}

type DockerContainer struct {
	ID         string
	Host       string
	Name       string
	HostIP     string
	HostPort   string
	TCP        string
	Config     *container.Config
	CLI        *docker.Client
	CustomArgs map[string]interface{}
}

func (d *DockerContainer) SetContainerID(id string) {
	d.ID = id
}

func (d *DockerContainer) SetContainerName(name string) {
	d.Name = name
}

func (d *DockerContainer) SetContainerHost(host string) {
	d.Host = host
}

func (d *DockerContainer) SetContainerHostIP(hostIp string) {
	d.HostIP = hostIp
}

func (d *DockerContainer) SetContainerHostPort(hostPort string) {
	d.HostPort = hostPort
}

func (d *DockerContainer) SetTCP(tcp string) {
	d.TCP = tcp
}

func (d *DockerContainer) SetConfig(config *container.Config) {
	d.Config = config
}

func (d *DockerContainer) SetContainerCLI(client *docker.Client) {
	d.CLI = client
}

func (d *DockerContainer) CreateContainer() *DockerContainer {
	ctx := context.Background()
	hostBinding := nat.PortBinding{
		HostIP:   d.HostIP,
		HostPort: d.HostPort,
	}
	containerPort, err := nat.NewPort("tcp", d.TCP)

	if err != nil {
		panic(err)
	}

	portBinding := nat.PortMap{containerPort: []nat.PortBinding{hostBinding}}
	hcfg := &container.HostConfig{
		PortBindings: portBinding,
	}
	cnt, err := d.CLI.ContainerCreate(ctx, d.Config, hcfg, nil, d.Name)

	if err != nil {
		panic(err)
	}

	if err = d.CLI.ContainerStart(ctx, cnt.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	d.SetContainerID(cnt.ID)

	return d
}
