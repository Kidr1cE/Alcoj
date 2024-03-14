package main

import (
	"alcoj/pkg/worker"
)

func main() {
	stop := make(chan struct{})
	go func() {
		worker.Run(stop)
		defer close(stop)
	}()
	<-stop
}
