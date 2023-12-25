package docker

import (
	"archive/tar"
	"bytes"
	"context"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

const (
	AppFolderPath  = "/app/source"
	DockerfilePath = "/app/dockerfile/Dockerfile"
)

type DockerInterface interface {
	Build()
	Create()
	Run()
}

type Environment struct {
	Raw        bool
	ImageName  string
	Dockerfile string
}

type Code struct {
	Type      int // 0: single file, 1: tar.gz, 2: zip
	Source    []byte
	TimeLimit int
}

type DockerWorker struct {
	ID    string
	Image string
	cli   *client.Client
}

func NewWorker() (*DockerWorker, error) {
	worker := new(DockerWorker)

	// docker client init
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	worker.cli = cli
	return worker, nil
}

func (d *DockerWorker) Info(ctx context.Context) (string, error) {
	cli := d.cli
	info, err := cli.Info(ctx)
	if err != nil {
		return "", err
	}
	return info.ID, nil
}

// Build docker image
func (d *DockerWorker) Build(ctx context.Context, dockerContext io.Reader) error {
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

	return nil
}

// Create docker container
func (d *DockerWorker) Create(ctx context.Context) error {
	resp, err := d.cli.ContainerCreate(ctx, &container.Config{
		Image:      d.Image,
		Entrypoint: []string{"bash", "/app/run.sh"},
	}, &container.HostConfig{
		Binds: []string{
			AppFolderPath,
		},
	}, nil, nil, "sandbox")
	if err != nil {
		return err
	}
	d.ID = resp.ID
	return nil
}

// Run docker container
func (d *DockerWorker) Run(ctx context.Context, input string) (string, string) {
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
		}
	case <-statusCh:
	}

	// get container logs/outputs
	out, err := cli.ContainerLogs(ctx, d.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		log.Fatalf("Container get logs failed: %v", err)
	}

	var codeOutputs bytes.Buffer

	_, err = stdcopy.StdCopy(&codeOutputs, &codeOutputs, out)
	// _, err = stdcopy.StdCopy(&buf, &buf, out)

	if err != nil {
		log.Fatalf("StdCopy failed: %v", err)
	}
	return codeOutputs.String(), ""
}

func (d *DockerWorker) Clean(ctx context.Context) error {
	cli := d.cli
	if err := cli.ContainerRemove(ctx, d.ID, types.ContainerRemoveOptions{}); err != nil {
		return err
	}
	return nil
}

// /app + dockerfile
func GetDockerContext(path string) (io.Reader, error) {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)

	// add Dockerfile
	dockerfilePath := filepath.Join(path, "dockerfile")
	if err := addFileToTarWriter(tw, dockerfilePath); err != nil {
		return nil, err
	}

	// add code dictionary
	appPath := filepath.Join(path, "app")
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