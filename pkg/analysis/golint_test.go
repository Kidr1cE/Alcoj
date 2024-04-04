package analysis

import (
	"alcoj/pkg/docker"
	"context"
	"fmt"
	"testing"
)

func TestGolint(t *testing.T) {
	containerID := "d7b2e6620765e6e169dd739684dc62e8d5cbb9252b641f07854f36096692e258"
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
	fmt.Println("=====================")
	message := parseGolintOutput(outputs, scriptPath)
	fmt.Println(message)
}
