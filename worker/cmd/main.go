package main

import (
	"context"
	"fmt"
	"worker/pkg/docker"
)

func main() {
	worker, err := docker.GetWorker()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	worker.Create(ctx)
	fmt.Println(worker.ID)
	fmt.Println(worker.Run(ctx, ""))
}
