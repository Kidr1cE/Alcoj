package main

import (
	"alcoj/pkg/worker"
)

func main() {
	stop := make(chan struct{})
	go worker.Run(stop)
	<-stop
}
