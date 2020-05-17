package mysql

import (
	"context"
	"docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/container"
	"encoding/json"
	"github.com/docker/go-connections/nat"
	"testing"
)

// Container ..
type Container struct {
	ID       string
	Host     string
	Name     string
	HostIP   string
	HostPort string
	TCP      string
	Config   *container.Config
	CLI      *docker.Client
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

// NewContainer creates a new container instance
func NewContainer(hp, tcp, n string, c *container.Config, cli *docker.Client) *Container {
	return &Container{
		Name:     n,
		HostPort: hp,
		TCP:      tcp,
		Config:   c,
		CLI:      cli,
	}
}

// StopContainer stops and removes the specified container.
func (c *Container) StopContainer() {
	ctx := context.Background()

	err := c.CLI.ContainerStop(ctx, c.ID, nil)

	if err != nil {
		panic(err)
	}
}

// DumpContainerLogs runs "docker logs" against the container and send it to t.Log
// TODO: tidy up
func (c *Container) DumpContainerLogs(t *testing.T) {
	ctx := context.Background()

	_, _ = c.CLI.ContainerLogs(ctx, c.ID, types.ContainerLogsOptions{})
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

func (c *Container) StartContainer() *Container {
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

	cnt, err := c.CLI.ContainerCreate(ctx, c.Config, hcfg, nil, "")

	if err != nil {
		panic(err)
	}

	if err = c.CLI.ContainerStart(ctx, cnt.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	c.ID = cnt.ID

	return c
}
