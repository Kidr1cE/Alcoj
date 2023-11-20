package docker

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// addFileToTarWriter 添加单个文件到 tar 归档
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

	// 重要：更改 header.Name 以反映文件在容器内的预期路径
	header.Name = strings.TrimPrefix(filename, "worker/docker/")
	if err := tw.WriteHeader(header); err != nil {
		return err
	}

	_, err = io.Copy(tw, file)
	return err
}

func createDockerContext() (io.Reader, error) {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)

	// 添加 Dockerfile
	if err := addFileToTarWriter(tw, "worker/docker/dockerfile"); err != nil {
		return nil, err
	}

	// 添加 code 目录下的所有文件
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

	if err := tw.Close(); err != nil {
		return nil, err
	}

	return buf, nil
}

func Run() {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	dockerBuildContext, err := createDockerContext()
	if err != nil {
		panic(err)
	}
	tags := fmt.Sprintf("%s:%s", "testimg", "0.0.1")

	buildOptions := types.ImageBuildOptions{
		Dockerfile: "dockerfile", // optional, is the default
		Tags:       []string{tags},
		Remove:     true,
	}

	output, err := cli.ImageBuild(context.Background(), dockerBuildContext, buildOptions)
	if err != nil {
		panic(err)
	}

	body, err := io.ReadAll(output.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Build resource image output: %v\n", string(body))

	if strings.Contains(string(body), "error") {
		panic("build image to docker error")
	}
}
