package docker

import (
	"archive/tar"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/google/uuid"
)

const (
	AppPath        = "/app"
	AppFolderPath  = "/app/source"
	DockerfilePath = "/app/dockerfile"
	EntryPath      = "/app/entry"
)

type Docker interface {
	Info(ctx context.Context) (string, error)
	Pull()
	Build()
	Create()
	Start()
	Run()
	Clean()
}

type DockerClient struct {
	ID         string
	Suffix     string
	Image      string
	Lang       string
	Version    string
	Status     int
	Raw        bool
	cli        *client.Client
	execConfig types.ExecConfig
}

func NewWorker() (*DockerClient, error) {
	// docker client init
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	worker := &DockerClient{
		ID:     "",
		Suffix: uuid.New().String(),
		Status: Unknown,
		cli:    cli,
	}

	return worker, nil
}

func (d *DockerClient) Info(ctx context.Context) (string, error) {
	cli := d.cli
	ctx, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()

	info, err := cli.Info(ctx)
	if err != nil {
		return "", err
	}

	return info.ID, nil
}

func (d *DockerClient) Pull(ctx context.Context) error {
	if d.Image == "" {
		return errors.New("image name is empty")
	}

	cli := d.cli
	reader, err := cli.ImagePull(ctx, d.Image, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer reader.Close()

	io.Copy(os.Stdout, reader)
	return nil
}

// Build docker image
func (d *DockerClient) Build(ctx context.Context) error {
	dockerContext, err := getDockerContext(AppPath)
	if err != nil {
		log.Fatalf("get docker context failed: %v", err)
		return err
	}

	// build image
	resp, err := d.cli.ImageBuild(ctx, dockerContext, types.ImageBuildOptions{
		Tags: []string{d.Image},
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// get build logs
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		return err
	}

	d.Status = NoSource

	return nil
}

// Create docker container
func (d *DockerClient) Create(ctx context.Context) error {
	var resp container.CreateResponse
	var err error
	if d.Raw {
		resp, err = d.cli.ContainerCreate(ctx, &container.Config{
			Image: d.Image,
		}, &container.HostConfig{
			Binds: []string{
				AppFolderPath + ":/app",
			},
		}, nil, nil, d.Lang)

		if err != nil {
			return err
		}
		d.execConfig = types.ExecConfig{
			AttachStdout: true,
			AttachStderr: true,
			Cmd:          []string{"bash", "-c", EntryPath + "/run.sh"},
		}
	} else {
		resp, err = d.cli.ContainerCreate(ctx, &container.Config{
			Image: d.Image,
			// Entrypoint: []string{"bash", "/app/run.sh"},
		}, &container.HostConfig{
			Binds: []string{
				AppFolderPath + ":/app",
			},
		}, nil, nil, d.Lang)
	}
	if err != nil {
		return err
	}

	d.ID = resp.ID
	return nil
}

// Run docker container
func (d *DockerClient) Start(ctx context.Context, input string) error {
	cli := d.cli

	if err := cli.ContainerStart(ctx, d.ID, types.ContainerStartOptions{}); err != nil {
		log.Fatalf("Container create failed: %v", err)
	}

	// wait
	statusCh, errCh := cli.ContainerWait(ctx, d.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			log.Fatalf("Container wait failed: %v", err)
			return err
		}
	case <-statusCh:
	}
	return nil
}

func (d *DockerClient) Run(ctx context.Context) error {
	cli := d.cli
	execID, err := cli.ContainerExecCreate(ctx, "container_name_or_id", d.execConfig)
	if err != nil {
		panic(err)
	}
	resp, err := cli.ContainerExecAttach(ctx, execID.ID, types.ExecStartCheck{})
	if err != nil {
		panic(err)
	}
	defer resp.Close()
	// read outputs

	return nil
}

func (d *DockerClient) Clean(ctx context.Context) error {
	cli := d.cli
	if err := cli.ContainerRemove(ctx, d.ID, types.ContainerRemoveOptions{}); err != nil {
		return err
	}
	return nil
}

// /app + dockerfile
func getDockerContext(path string) (io.Reader, error) {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)

	// add Dockerfile
	dockerfilePath := filepath.Join(path, "dockerfile/dockerfile")
	if err := addFileToTarWriter(tw, dockerfilePath); err != nil {
		return nil, err
	}

	// add code dictionary
	appPath := filepath.Join(path, "source")
	err := filepath.Walk(appPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return addFileToTarWriter(tw, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func addFileToTarWriter(tw *tar.Writer, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := tar.FileInfoHeader(stat, stat.Name())
	if err != nil {
		return err
	}

	header.Name = strings.TrimPrefix(filename, "pkg/docker/")
	if err := tw.WriteHeader(header); err != nil {
		return err
	}

	_, err = io.Copy(tw, file)
	return err
}
