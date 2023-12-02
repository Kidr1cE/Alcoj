package main

import (
	"context"
	"worker/pkg/docker"
)

func main() {
	worker, err := docker.GetWorker()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	worker.Create(ctx)
	worker.Run(ctx, "")
}
