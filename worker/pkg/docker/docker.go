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
	"worker/pkg/utils"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

type Docker interface {
	Build()
	Create()
	Run()
}
type DockerConfig struct {
	Image string
	Cmd   string
}
type DockerWorker struct {
	ID  string
	cli *client.Client
}

func GetWorker() (*DockerWorker, error) {
	worker := new(DockerWorker)

	// docker client init
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	worker.cli = cli
	return worker, nil
}

func (d *DockerWorker) Build(ctx context.Context) error {

	return nil
}

func (d *DockerWorker) Create(ctx context.Context) error {
	path, err := utils.GetPath("pkg/docker/app")
	if err != nil {
		return err
	}

	resp, err := d.cli.ContainerCreate(ctx, &container.Config{
		Image:      "golang:1.20.11",
		Entrypoint: []string{"bash", "/app/run.sh"},
	}, &container.HostConfig{
		Binds: []string{
			path + ":/app",
		},
	}, nil, nil, "sandbox")
	if err != nil {
		return err
	}
	d.ID = resp.ID
	return nil
}

func (d *DockerWorker) Run(ctx context.Context, input string) string {
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

	var buf bytes.Buffer
	_, err = stdcopy.StdCopy(&buf, &buf, out)
	if err != nil {
		log.Fatalf("StdCopy failed: %v", err)
	}
	return buf.String()
}

func getDockerContext() (io.Reader, error) {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)

	// add Dockerfile
	if err := addFileToTarWriter(tw, "worker/docker/dockerfile"); err != nil {
		return nil, err
	}

	// add code dictionary
	err := filepath.Walk("worker/docker/code", func(path string, info os.FileInfo, err error) error {
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

	header.Name = strings.TrimPrefix(filename, "worker/docker/")
	if err := tw.WriteHeader(header); err != nil {
		return err
	}

	_, err = io.Copy(tw, file)
	return err
}
