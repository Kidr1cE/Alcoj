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

func DockerPoolInit() {
	ctx := context.Background()
	// docker client init
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "golang:1.20.11",
		Cmd:   []string{"go", "run", "/app/code/main.go"},
	}, &container.HostConfig{
		Binds: []string{
			"/home/alco/go-project/golang-oj-worker/worker/docker/code:/app/code",
		},
	}, nil, nil, "golang-runner")
	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		log.Fatalf("Container create failed: %v", err)
	}

	// wait
	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			log.Fatalf("Container wait failed: %v", err)
		}
	case <-statusCh:
	}

	// get container logs/outputs
	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		log.Fatalf("Container get logs failed: %v", err)
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)
}
