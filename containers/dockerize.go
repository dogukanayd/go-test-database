package containers

import (
	"context"
	"docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/container"
	"encoding/json"
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
	CreateContainer() *Container
	StopContainer()
}

// DockerInspect Values resulting from the command `docker inspect container-name`
type DockerInspect struct {
	State struct {
		Status string `json:"Status"`
	} `json:"State"`
	NetworkSettings struct {
		Networks struct {
			Bridge struct {
				IPAddress string `json:"IPAddress"`
			} `json:"bridge"`
		} `json:"Networks"`
	} `json:"NetworkSettings"`
}

type Container struct {
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

func (c *Container) SetContainerID(id string) {
	c.ID = id
}

func (c *Container) SetContainerName(name string) {
	c.Name = name
}

func (c *Container) SetContainerHost(host string) {
	c.Host = host
}

func (c *Container) SetContainerHostIP(hostIp string) {
	c.HostIP = hostIp
}

func (c *Container) SetContainerHostPort(hostPort string) {
	c.HostPort = hostPort
}

func (c *Container) SetTCP(tcp string) {
	c.TCP = tcp
}

func (c *Container) SetConfig(config *container.Config) {
	c.Config = config
}

func (c *Container) SetContainerCLI(client *docker.Client) {
	c.CLI = client
}

func (c *Container) CreateContainer() *Container {
	ctx := context.Background()
	hostBinding := nat.PortBinding{
		HostIP:   c.HostIP,
		HostPort: c.HostPort,
	}
	containerPort, err := nat.NewPort("tcp", c.TCP)

	if err != nil {
		panic(err)
	}

	portBinding := nat.PortMap{containerPort: []nat.PortBinding{hostBinding}}
	hcfg := &container.HostConfig{
		PortBindings: portBinding,
	}
	cnt, err := c.CLI.ContainerCreate(ctx, c.Config, hcfg, nil, c.Name)

	if err != nil {
		panic(err)
	}

	if err = c.CLI.ContainerStart(ctx, cnt.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	c.SetContainerID(cnt.ID)

	return c
}

func (c *Container) GetStatus() ([]DockerInspect, error) {
	ctx := context.Background()
	var di []DockerInspect

	_, by, _ := c.CLI.ContainerInspectWithRaw(ctx, c.ID, true)

	if err := json.Unmarshal(by, &di); err != nil {
		return di, err
	}

	return di, nil
}

// StopContainer ...
func (c *Container) StopContainer() {
	ctx := context.Background()

	err := c.CLI.ContainerStop(ctx, c.ID, nil)

	if err != nil {
		panic(err)
	}
}

// DumpContainerLogs runs "docker logs" against the container and send it to t.Log
// TODO: tidy up
func (c *Container) DumpContainerLogs() {
	ctx := context.Background()

	_, _ = c.CLI.ContainerLogs(ctx, c.ID, types.ContainerLogsOptions{})
}
