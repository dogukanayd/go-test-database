package testmysql

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"testing"
	"time"
)

// Container ..
type Container struct {
	ID   string
	Host string
	Name string
}

const ContainerName = "mysql_test"

// StartContainer ..
func StartContainer(t *testing.T) *Container {
	t.Helper()

	startContainerAndMigrate(t)

	cmd := exec.Command("docker", "ps", "-aqf", "name=" + ContainerName)

	var out bytes.Buffer
	var stderr bytes.Buffer

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
		t.Fatalf("could not inspect container %s: %v", ci, stderr.String())

		return nil
	}

	var doc []struct {
		NetworkSettings struct {
			Ports struct {
				TCP3306 []struct {
					HostIP   string `json:"HostIp"`
					HostPort string `json:"HostPort"`
				} `json:"3306/tcp"`
			} `json:"Ports"`
		} `json:"NetworkSettings"`
	}

	if err := json.Unmarshal(out.Bytes(), &doc); err != nil {
		t.Fatalf("could not decode json: %v", err)
	}

	network := doc[0].NetworkSettings.Ports.TCP3306[0]

	c := Container{
		ID:   ci,
		Host: network.HostIP + ":" + network.HostPort,
		Name: ContainerName,
	}

	t.Log("DB Host: ", c.Host)
	t.Log("DB Container: ", c.ID)

	return &c
}

func startContainerAndMigrate(t *testing.T)  {
	var out bytes.Buffer
	var stderr bytes.Buffer
	t.Helper()

	cmd := exec.Command("bash", "run.sh")
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		t.Fatalf("could not build container %v", stderr.String())
	}

	maxAttempts := 20

	for attempts := 1; attempts <= maxAttempts; attempts++ {
		if out.String() == "completed" {
			t.Log("test database up success")
		}

		time.Sleep(time.Duration(attempts) * 100 * time.Millisecond)
	}
}

// StopContainer stops and removes the specified container.
// TODO: it should stop container by id
func StopContainer(t *testing.T, c *Container) {
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
func DumpContainerLogs(t *testing.T, c *Container) {
	t.Helper()

	out, err := exec.Command("docker", "logs", c.ID).CombinedOutput()

	if err != nil {
		t.Fatalf("could not log container: %v", err)
	}

	t.Logf("Logs for %s\n%s:", c.ID, out)
}
