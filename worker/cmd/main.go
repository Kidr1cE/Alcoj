package main

import "worker/pkg/docker"

func main() {
	worker, err := docker.GetWorker()
	if err != nil {
		panic(err)
	}
	worker.RunCode("")
}
