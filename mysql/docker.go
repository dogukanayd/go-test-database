package mysql

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"testing"
)

const ContainerName = "mysql_test"

// Container ..
type Container struct {
	ID   string
	Host string
	Name string
}

// DockerInspect Values resulting from the command `docker inspect container-name`
type DockerInspect struct {
	NetworkSettings struct {
		Ports struct {
			TCP3306 []struct {
				HostIP   string `json:"HostIp"`
				HostPort string `json:"HostPort"`
			} `json:"3306/tcp"`
		} `json:"Ports"`
	} `json:"NetworkSettings"`
}

// NewContainer creates a new container instance
func NewContainer() *Container {
	var doc []DockerInspect
	var out bytes.Buffer
	var stderr bytes.Buffer
	var c Container

	c.StartContainer()

	cmd := exec.Command("docker", "ps", "-aqf", "name="+ContainerName)
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	_ = cmd.Run()

	ci := out.String()

	cmd = exec.Command("docker", "inspect", ContainerName)
	out.Reset()
	stderr.Reset()
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		fmt.Errorf("could not inspect container %s: %v", ci, stderr.String())

		return nil
	}

	if err := json.Unmarshal(out.Bytes(), &doc); err != nil {
		fmt.Errorf("could not decode docker inspect data: %v", err)
	}

	network := doc[0].NetworkSettings.Ports.TCP3306[0]
	c.ID = ci
	c.Host = network.HostIP + ":" + network.HostPort
	c.Name = ContainerName

	return &c
}

// StopContainer stops and removes the specified container.
// TODO: it should stop container by id, right now its giving an error when trying to use c.ID figure it out why?
func (c *Container) StopContainer(t *testing.T) {
	t.Helper()

	t.Log("container id", c.ID)

	if err := exec.Command("docker", "kill", ContainerName).Run(); err != nil {
		t.Fatalf("could not stop container: %v", err)
	}

	t.Log("Stopped:", c.ID)

	if err := exec.Command("docker", "container", "rm", "-f", ContainerName).Run(); err != nil {
		t.Fatalf("could not remove container: %v", err)
	}

	t.Log("Removed:", c.ID)
}

// DumpContainerLogs runs "docker logs" against the container and send it to t.Log
func (c *Container) DumpContainerLogs(t *testing.T) {
	t.Helper()

	out, err := exec.Command("docker", "logs", c.ID).CombinedOutput()

	if err != nil {
		t.Fatalf("could not log container: %v", err)
	}

	t.Logf("Logs for %s\n%s:", c.ID, out)
}

func (c *Container) StartContainer() {
	killAndRemove()
	runContainer()
}

func killAndRemove() {
	_ = exec.Command("docker", "kill", "mysql_test").Run()
	_ = exec.Command("docker", "container", "rm", "-f", "mysql_test").Run()
}

func runContainer() {
	_ = exec.Command("docker", "run", "-d", "--name", "mysql_test", "-p", "3305:3306", "-e", "MYSQL_ROOT_PASSWORD=root", "-e", "MYSQL_DATABASE=test_database", "mysql").Run()
	_ = exec.Command("docker", "wait", "mysql_test")
}
