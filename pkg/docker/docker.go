package docker

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

const (
	Unknown = iota
	Ready
	Running
)

type DockerClient struct {
	ID          string
	Image       string
	ContainerID string
	Status      int
	cli         *client.Client
}

func NewDocker(id string) (*DockerClient, error) {
	// docker client init
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	worker := &DockerClient{
		ID:     id,
		Status: Unknown,
		cli:    cli,
	}

	return worker, nil
}

// Check if docker is installed
func (d *DockerClient) Info(ctx context.Context) (string, error) {
	cli := d.cli
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	info, err := cli.Info(ctx)
	if err != nil {
		return "", err
	}

	return info.ID, nil
}

func (d *DockerClient) CheckContainerHealth(ctx context.Context) (bool, error) {
	cli := d.cli
	id := d.ContainerID

	// check container status
	resp, err := cli.ContainerInspect(ctx, id)
	if err != nil {
		if client.IsErrNotFound(err) {
			log.Println("Container not found")
			return false, nil
		} else {
			return false, err
		}
	}

	if resp.State.Running {
		return true, nil
	} else {
		return false, nil
	}
}

// Create docker container
func (d *DockerClient) Create(ctx context.Context) error {
	var resp container.CreateResponse

	resp, err := d.cli.ContainerCreate(ctx, &container.Config{
		Image:      d.Image,
		WorkingDir: "/sandbox",
		Cmd:        []string{"tail", "-f", "/dev/null"},
	}, &container.HostConfig{
		Binds: []string{
			"sandbox:/sandbox",
		},
	}, nil, nil, "containerName")

	if err != nil {
		return err
	}

	d.ContainerID = resp.ID
	return nil
}

// Run docker container
func (d *DockerClient) Start(ctx context.Context, input string) error {
	cli := d.cli

	if err := cli.ContainerStart(ctx, d.ContainerID, types.ContainerStartOptions{}); err != nil {
		log.Fatalf("Container create failed: %v", err)
	}

	return nil
}

func (d *DockerClient) Run(ctx context.Context) (string, error) {
	execConfig := types.ExecConfig{
		Cmd:          []string{"bash", "/sandbox" + "/run.sh"},
		AttachStdout: true,
		AttachStderr: true,
	}

	execID, err := d.cli.ContainerExecCreate(ctx, d.ID, execConfig)
	if err != nil {
		return "", err
	}

	resp, err := d.cli.ContainerExecAttach(ctx, execID.ID, types.ExecStartCheck{})
	if err != nil {
		return "", err
	}
	defer resp.Close()

	output, err := io.ReadAll(resp.Reader)
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func (d *DockerClient) Cmd(ctx context.Context, cmd []string) (string, string, error) {
	execConfig := types.ExecConfig{
		Cmd:          cmd,
		AttachStdout: true,
		AttachStderr: true,
	}

	execID, err := d.cli.ContainerExecCreate(ctx, d.ID, execConfig)
	if err != nil {
		return "", "", fmt.Errorf("failed to create exec instance in container: %w", err)
	}

	resp, err := d.cli.ContainerExecAttach(ctx, execID.ID, types.ExecStartCheck{})
	if err != nil {
		return "", "", fmt.Errorf("failed to attach to exec instance: %w", err)
	}
	defer resp.Close()

	return "", "", nil
}

func (d *DockerClient) Clean(ctx context.Context) error {
	cli := d.cli
	if err := cli.ContainerRemove(ctx, d.ID, types.ContainerRemoveOptions{}); err != nil {
		return err
	}
	return nil
}
