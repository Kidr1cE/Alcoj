package main

import (
	"worker/docker"
)

func main() {
	worker, err := docker.GetWorker()
	if err != nil {
		panic(err)
	}
	worker.RunCode("")
}
