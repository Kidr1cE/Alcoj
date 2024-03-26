package analysis

import (
	"alcoj/pkg/docker"
	"context"
	"fmt"
	"testing"
)

func TestGolint(t *testing.T) {
	containerID := "8b0e672c0aff3f27e3a95aeb5c1056d9af9e29f913d191e56b515604f19d660e"
	scriptPath := "/sandbox/main.go"

	dockerClient, err := docker.NewDocker("some-id")
	if err != nil {
		t.Errorf("NewDocker() failed: %v", err)
	}
	dockerClient.ContainerID = containerID

	outputs, err := dockerClient.Cmd(context.Background(), []string{"golangci-lint", "run", scriptPath}, "")
	if err != nil {
		t.Errorf("Cmd() failed: %v", err)
	}

	fmt.Println("\n=====================")
	fmt.Println(outputs)
}
