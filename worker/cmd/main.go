package main

import (
	"alcoj/worker/pkg/docker"
	"context"
	"fmt"
)

var worker *docker.DockerWorker

func main() {
	worker, err := docker.GetWorker(&docker.Environment{Code: "golang", Version: "1.20.11"})
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	worker.Create(ctx)
	fmt.Println(worker.ID)
	fmt.Println(worker.Run(ctx, ""))
}
